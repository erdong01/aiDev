curl "http://127.0.0.1:9180/apisix/admin/routes/ai-proxy" \
-X PUT \
-H "X-API-KEY: FBB7C3831B96B56ECF00052E4B773D8A" \
-d '{
  "uri": "/v1/pqcontents/gen/tasks",
  "methods": ["POST"],
  "upstream_id": "volcengine-ark",
  "plugins": {
    "key-auth": {},
    "proxy-rewrite": { "uri": "/api/v3/contents/generations/tasks" },
    "response-rewrite": {}, 
    "serverless-post-function": {
      "functions": [
        "return function(conf, ctx)
            local core = require(\"apisix.core\")
            local http = require(\"resty.http\")
            
            -- 核心修复：确保能拿到响应体
            local resp_body = ctx.resp_body
            if not resp_body then 
                core.log.warn(\"警告：未能捕获到上游响应体，回调终止\")
                return 
            end

            local req_body = core.request.get_body() or \"\"

            -- 获取已经注入的火山 Key
            local volc_key = ngx.var.http_authorization
            local consumer_name = ctx.consumer_name or \"unknown\"

            local function sync_task_to_go()
                local httpc = http.new()
                httpc:set_timeout(20000) -- 5秒超时
                local callback_body = core.json.encode({
                    request = {
                        headers = {
                            authorization = volc_key
                        },
                        body = req_body
                    },
                    response = {
                        status = ngx.status,
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
  "desc": "开启响应体缓存并回调 Go 入库"
}'
