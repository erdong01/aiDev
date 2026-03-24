docker build --platform linux/amd64 -t registry.cn-hangzhou.aliyuncs.com/panqu/aidev:1.0.11 .

docker push registry.cn-hangzhou.aliyuncs.com/panqu/aidev:1.0.12

docker pull registry.cn-hangzhou.aliyuncs.com/panqu/aidev:1.0.12

docker stop aidev

docker rm aidev
