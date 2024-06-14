# web

## 说明

web是jrasp的管理端工程，负责配置策略和日志消费。前端由vue开发、后端使用golang。

## 预览

https://www.server.jrasp.com:8088/rasp-admin
账号和密码： admin 123456

## 开发/编译

由于本产品安装数量较多，基于产品安全考虑，不再公开源码，合作伙伴可提供jrasp基础框架源码

### 运行

解压之后执行（内置数据库、前后端不分离）
```
./server
```
```
MacBook-Pro root$ ./server 
2023-05-01 13:48:57     INFO    common/logger.go:108    server/common.InitLogger        初始化zap日志完成!
2023-05-01 13:48:58     INFO    common/database.go:30   server/common.InitMysql 初始化数据库完成!
2023-05-01 13:48:58     INFO    common/casbin.go:22     server/common.InitCasbinEnforcer        初始化Casbin完成!
2023-05-01 13:48:58     INFO    common/validator.go:26  server/common.InitValidate      初始化validator.v10数据校验器完成
2023-05-01 13:48:58     INFO    routes/routes.go:73     server/routes.InitRoutes        初始化路由完成！
2023-05-01 13:48:58     INFO    server/main.go:84       main.main       Server is running at 0.0.0.0:8088/rasp-admin
2023-05-01 13:48:58     INFO    socket/ws.go:27 server/socket.(*Manager).Start  websocket manage start
```

以后台进程启动
```shell
nohup ./server >/dev/null 2>&1 &
```

## 访问

+ 访问 `http://localhost:8088/rasp-admin` 
+ 账号和密码： admin 123456

## 替换证书（非必须）
默认启动方式为非http方式，可以更新ssl证书并开启ssl功能。修改`config.yml`中的如下配置
```yaml
ssl:
  # https开关,默认关闭
  enable: false
  # ssl 证书key
  keyFile: keyFile.key
  # ssl 证书路径
  certFile: certFile.pem
```

## 使用mysql数据库（非必须）
默认使用sqlite数据库，可以替换为mysql数据库。修改`config.yml`中的如下配置
```yaml
database:
  # driver: mysql, sqlite
  driver: sqlite
#  source: root:Gepoint@tcp(127.0.0.1:3306)/jrasp?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&timeout=1000ms
  source: jrasp.db
  # silent=1 error=2 warn=3 info=4
  log-mode: 1
```
