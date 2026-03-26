curl "http://127.0.0.1:9180/apisix/admin/routes/ai-proxy" \
-X PUT \
-H "X-API-KEY: FBB7C3831B96B56ECF00052E4B773D8A" \
-d '{
  "uri": "/v1/pqcontents/gen/tasks",
  "methods": ["POST"],
  "upstream_id": "volcengine-ark",
  "plugins": {
    "key-auth": {},
    "proxy-rewrite": {
      "uri": "/api/v3/contents/generations/tasks",
      "headers": {
        "set": {
          "Accept-Encoding": "identity"
        }
      }
    },
    "serverless-pre-function": {
      "phase": "rewrite",
      "functions": [
        "return function(conf, ctx)
            local core = require(\"apisix.core\")
            local function preview_payload(value, limit)
                if not value then
                    return \"nil\"
                end

                local preview = value
                if #preview > limit then
                    preview = string.sub(preview, 1, limit) .. \"...(truncated)\"
                end

                preview = string.gsub(preview, \"[\\r\\n\\t]\", \" \")
                return preview
            end

            local headers = ngx.req.get_headers()
            core.log.warn(
                \"[aidev-video-debug] rewrite 进入, content_type=\",
                headers[\"content-type\"] or \"nil\",
                \", content_length=\",
                headers[\"content-length\"] or \"nil\",
                \", transfer_encoding=\",
                headers[\"transfer-encoding\"] or \"nil\"
            )

            -- 在请求阶段一次性读取并缓存原始 body，后续阶段不再调用 req API。
            ngx.req.read_body()

            local req_body = ngx.req.get_body_data()
            if not req_body then
                local body_file = ngx.req.get_body_file()
                if body_file then
                    core.log.warn(\"[aidev-video-debug] rewrite 请求体落盘, body_file=\", body_file)
                    local file = io.open(body_file, \"rb\")
                    if file then
                        req_body = file:read(\"*a\")
                        file:close()
                    end
                end
            end

            if not req_body or req_body == \"\" then
                core.log.warn(\"[aidev-video-debug] rewrite 阶段未拿到原始请求体，回退为 {}\")
                req_body = \"{}\"
            end

            ctx.cached_req_body = ngx.encode_base64(req_body)
            ctx.cached_req_body_base64 = true

            core.log.warn(
                \"[aidev-video-debug] rewrite 请求体已缓存, len=\",
                #req_body,
                \", cached_len=\",
                #ctx.cached_req_body,
                \", req_body_base64=\",
                tostring(ctx.cached_req_body_base64),
                \", preview=\",
                preview_payload(req_body, 600)
            )

            ctx.cached_authorization = ngx.var.http_authorization
            ctx.cached_consumer_name = ctx.consumer_name or \"unknown\"
        end"
      ]
    },
    "serverless-post-function": {
      "phase": "body_filter",
      "functions": [
        "return function(conf, ctx)
            local core = require(\"apisix.core\")
            local http = require(\"resty.http\")
            local function preview_payload(value, limit)
                if not value then
                    return \"nil\"
                end

                local preview = value
                if #preview > limit then
                    preview = string.sub(preview, 1, limit) .. \"...(truncated)\"
                end

                preview = string.gsub(preview, \"[\\r\\n\\t]\", \" \")
                return preview
            end

            local chunk = ngx.arg[1]
            local eof = ngx.arg[2]

            ctx.resp_body = (ctx.resp_body or \"\") .. (chunk or \"\")

            -- body_filter 会执行多次，只在最后一个 chunk 到来时回调一次。
            if not eof then
                return
            end

            local resp_body = ctx.resp_body
            if not resp_body or resp_body == \"\" then
                core.log.warn(\"警告：未能捕获到上游响应体，回调终止\")
                return
            end

            local req_body = ctx.cached_req_body or ngx.encode_base64(\"{}\")
            local req_body_base64 = ctx.cached_req_body_base64 == true

            local volc_key = ctx.cached_authorization
            local consumer_name = ctx.cached_consumer_name or \"unknown\"
            local response_status = ngx.status

            core.log.warn(
                \"[aidev-video-debug] body_filter 准备回调, status=\",
                response_status,
                \", req_body_base64=\",
                tostring(req_body_base64),
                \", req_body_len=\",
                #req_body,
                \", req_body_preview=\",
                preview_payload(req_body, 240),
                \", resp_body_len=\",
                #resp_body
            )

            local function sync_task_to_go()
                local httpc = http.new()
                httpc:set_timeout(20000)
                local callback_body = core.json.encode({
                    request = {
                        headers = {
                            authorization = volc_key
                        },
                        body = req_body,
                        body_base64 = req_body_base64
                    },
                    response = {
                        status = response_status,
                        body = ngx.encode_base64(resp_body),
                        body_base64 = true
                    },
                    consumer = {
                        username = consumer_name
                    }
                })

                local res, err = httpc:request_uri(\"http://aidev:8080/v1/volcengine/callback\", {
                    method = \"POST\",
                    body = callback_body,
                    headers = {
                        [\"Content-Type\"] = \"application/json\",
                        [\"Authorization\"] = volc_key,
                        [\"X-Consumer-Name\"] = consumer_name
                    }
                })
                if not res then
                    core.log.error(\"回调 Go 服务失败: \", err)
                else
                    core.log.info(\"回调 Go 服务成功: \", res.status)
                end
            end

            ngx.timer.at(0, sync_task_to_go)
        end"
      ]
    }
  },
  "desc": "在 body_filter 阶段拼接响应体并回调 Go 入库"
}'
