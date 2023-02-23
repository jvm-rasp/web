# web

npm 版本 >= 8.5.0

golang 版本 >= 1.19.6 (强制)

## server 工程

进入到 server 目录下

### 下载依赖
```
export GOPROXY="https://mirrors.aliyun.com/goproxy/"
go mod tidy
```
### 编译
```
go build
```

### 运行
```
./server
```
```
2023-02-22 23:30:27     INFO    routes/routes.go:54     server/routes.InitRoutes        初始化路由完成！
2023-02-22 23:30:27     INFO    server/main.go:64       main.main       Server is running at localhost:8088/api
```


## UI 工程

进入到ui目录下

###  编译
首次编译时，更换阿里的 npm registry 镜像
```
npm config set registry https://registry.npm.taobao.org
```
### 安装依赖
```
npm install
```

### 本地运行
```
npm run dev
```

### 静态文件
```
npm run build:prod
```
## 登陆账号

admin
123456