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
      "phase": "access",
      "functions": [
        "return function(conf, ctx)
            local core = require(\"apisix.core\")

            -- 在请求阶段把后面要用的数据缓存到 ctx，避免在响应阶段拿不到。
            ctx.cached_req_body = core.request.get_body() or \"\"
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

            local req_body = ctx.cached_req_body or \"\"
            local volc_key = ctx.cached_authorization
            local consumer_name = ctx.cached_consumer_name or \"unknown\"
            local response_status = ngx.status

            local function sync_task_to_go()
                local httpc = http.new()
                httpc:set_timeout(20000)
                local callback_body = core.json.encode({
                    request = {
                        headers = {
                            authorization = volc_key
                        },
                        body = req_body
                    },
                    response = {
                        status = response_status,
                        body = resp_body
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
