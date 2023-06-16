# 配置https

## 证书替换

管理端配置https。默认提供了绑定www.server.jrasp.com的ssl证书。 证书替换可以联系patton解决。

## 腾讯云ssl证书配置

![ssl证书申请](./images/ssl-download.png)

![ssl证书下载](./images/ssl.png)

## 证书转换

```
pfx=www.server.jrasp.com.pfx
openssl pkcs12 -in $pfx -nocerts -out key.pem -nodes
openssl pkcs12 -in $pfx -nokeys -out test.pem
openssl rsa -in key.pem -out test.key
mv test.key  keyFile.key
mv test.pem  certFile.pem
```
复制到安装路径下

## server开启ssl
修改 config.yml ssl配置选项
```
ssl:
  # https开关,默认关闭
  enable: true
  # ssl 证书key
  keyFile: keyFile.key
  # ssl 证书路径
  certFile: certFile.pem
```
