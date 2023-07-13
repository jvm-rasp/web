#!/bin/bash
set -e

version=1.1.2
echo "build web project(Golang) start:" $(date +"%Y-%m-%d %H:%M:%S")

# 切换到 build.sh 脚本所在的目录
cd $(dirname $0)

BUILD_TMP=../target
if [ -d ${BUILD_TMP} ]; then
  rm -rf ${BUILD_TMP}
fi

# server target dir
BUILD_TARGET_DIR=../target/server

SERVER_PROJECT_HOME=../server

# ui
UI_PROJECT_HOME=../ui
cd ${UI_PROJECT_HOME}

npm config set registry https://registry.npm.taobao.org

npm install && npm run build:prod

if [ -d ${SERVER_PROJECT_HOME}/resources/html ]; then
  rm -rf ${SERVER_PROJECT_HOME}/resources/html
fi
mv ${UI_PROJECT_HOME}/dist  ${UI_PROJECT_HOME}/html
mv ${UI_PROJECT_HOME}/html  ${SERVER_PROJECT_HOME}/resources/

echo "build web project(Golang) start:" $(date +"%Y-%m-%d %H:%M:%S")


cd ${SERVER_PROJECT_HOME}
projectpath=`pwd`
echo "current dir:${projectpath}"

# 设置阿里云代理
#export GOPROXY="https://mirrors.aliyun.com/goproxy/"
export GOPROXY="goproxy.cn"
echo "GOPROXY:"${GOPROXY}

moduleName=$(go list -m)
program=$(basename ${moduleName})

cd ${projectpath}

go mod tidy

# sqlite 不支持交叉编译
# Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work
go build -v -ldflags "$flags" -o ${projectpath}/$program

echo "build go project end:" $(date +"%Y-%m-%d %H:%M:%S")

# reset the target dir
mkdir -p ${BUILD_TARGET_DIR}

# 文件复制
cp ${projectpath}/server ${BUILD_TARGET_DIR}/ && \
cp ${projectpath}/rbac_model.conf ${BUILD_TARGET_DIR}/ && \
cp ${projectpath}/web-priv.pem ${BUILD_TARGET_DIR}/ && \
cp ${projectpath}/web-pub.pem ${BUILD_TARGET_DIR}/ && \
cp ${projectpath}/keyFile.key ${BUILD_TARGET_DIR}/ && \
cp ${projectpath}/certFile.pem ${BUILD_TARGET_DIR}/ && \
cp ${projectpath}/config.yml ${BUILD_TARGET_DIR}/

# 打包

# make it execute able
chmod +x ${BUILD_TARGET_DIR}/server

os_name=$(uname -s | tr '[:upper:]' '[:lower:]')
arch_name=$(uname -m | tr '[:upper:]' '[:lower:]')

# tar the jrasp-server.tar.gz
cd ${BUILD_TMP} || exit_on_err 1 "[ERROR] cd target dir error."

zip_name="jrasp-server-${version}-${os_name}-${arch_name}-$(date +"%Y%m%d%H%M%S").tar.gz"

tar -zcvf ../$zip_name server/

echo "$zip_name finish."

exit 0