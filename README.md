docker build --platform linux/amd64 -t registry.cn-hangzhou.aliyuncs.com/panqu/aidev:1.0.20 .

docker push registry.cn-hangzhou.aliyuncs.com/panqu/aidev:1.0.20

docker pull registry.cn-hangzhou.aliyuncs.com/panqu/aidev:1.0.20

docker stop aidev

docker rm aidev


docker  logs ai_build-apisix-1 --tail 200