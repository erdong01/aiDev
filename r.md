docker build --platform linux/amd64 -t registry.cn-hangzhou.aliyuncs.com/panqu/aidev:1.0.11 .

docker push registry.cn-hangzhou.aliyuncs.com/panqu/aidev:1.0.11

docker pull registry.cn-hangzhou.aliyuncs.com/panqu/aidev:1.0.11

docker stop aidev

docker rm aidev

docker run -d -p 8080:8080 \
  -e MYSQL_USER=root \
  -e MYSQL_PASS='panQu!142857' \
  -e MYSQL_HOST=115.191.19.88 \
  --network ai_build_apisix \
  --name aidev \
  --restart always \
  registry.cn-hangzhou.aliyuncs.com/panqu/aidev:1.0.2