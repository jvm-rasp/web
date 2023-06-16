# web

## 说明

web是jrasp的管理端工程，负责配置策略和日志消费。前端由vue开发、后端使用golang。

## 预览

https://www.server.jrasp.com:8088/rasp-admin

账号和密码： admin 123456

(服务器带宽低，打开较慢)

## 编译工具

npm 版本 >= 8.5.0

golang 版本 >= 1.19.6 (强制)

## 一键打包

进入到 build 目录下执行

+ linux
```
bash build.sh
```

+ windows
```
（待补充）
```

输出文件在`target`目录下


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

## 访问

+ 访问 `http://localhost:8088/rasp-admin` 

+ 账号和密码： admin 123456

(默认关闭关闭https)

## 替换证书（非必须）

参考：[https配置](./https.md)



