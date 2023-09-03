# Nacos安装教程

### 启动nacos，官网的那个Apple M1目前有问题，启动不了（不推荐安装 MAC专用）
docker run --name hello-nacos -e MODE=standalone -p 8848:8848 -d zill057/nacos-server-apple-silicon:2.0.3

### 启动nacos，官方版本（单机模式）推荐安装

1.Clone 项目

git clone https://github.com/nacos-group/nacos-docker.git

cd nacos-docker

2.单机模式 Derby

docker-compose -f example/standalone-derby.yaml up


### 界面
```shell
打开浏览器访问：http://127.0.0.1:8848/nacos/
账号：nacos
密码：nacos

包括注册中心和配置中心，可以直接编辑配置文件使用
```

### 服务注册
curl -X POST 'http://127.0.0.1:8848/nacos/v1/ns/instance?serviceName=nacos.naming.serviceName&ip=20.18.7.10&port=8080'

### 服务发现
curl -X GET 'http://127.0.0.1:8848/nacos/v1/ns/instance/list?serviceName=nacos.naming.serviceName'

### 发布配置
curl -X POST "http://127.0.0.1:8848/nacos/v1/cs/configs?dataId=nacos.cfg.dataId&group=test&content=HelloWorld"

### 获取配置
curl -X GET "http://127.0.0.1:8848/nacos/v1/cs/configs?dataId=nacos.cfg.dataId&group=test"

