#!/bin/bash

echo "build web project(Golang) start:" $(date +"%Y-%m-%d %H:%M:%S")

# 切换到 build.sh 脚本所在的目录
cd $(dirname $0) || exit 1

# exit shell with err_code
# $1 : err_code
# $2 : err_msg
exit_on_err()
{
    [[ ! -z "${2}" ]] && echo "${2}" 1>&2
    exit ${1}
}

BUILD_TMP=../target
if [ -d ${BUILD_TMP} ]; then
  rm -rf ${BUILD_TMP} || exit_on_err 1 "[ERROR] delete target dir error."
fi

# server target dir
BUILD_TARGET_DIR=../target/server

SERVER_PROJECT_HOME=../server

# ui
UI_PROJECT_HOME=../ui
cd ${UI_PROJECT_HOME} || exit 1

npm config set registry https://registry.npm.taobao.org

npm install && npm run build:prod || exit 1

if [ -d ${SERVER_PROJECT_HOME}/resources/html ]; then
  rm -rf ${SERVER_PROJECT_HOME}/resources/html
fi
mv ${UI_PROJECT_HOME}/dist  ${UI_PROJECT_HOME}/html || exit 1
mv ${UI_PROJECT_HOME}/html  ${SERVER_PROJECT_HOME}/resources/ || exit 1

echo "build web project(Golang) start:" $(date +"%Y-%m-%d %H:%M:%S")


cd ${SERVER_PROJECT_HOME} || exit 1
projectpath=`pwd`
echo "current dir:${projectpath}"

# 设置代理
#export GOPROXY="https://mirrors.aliyun.com/goproxy/"
export GOPROXY="goproxy.cn"
echo "GOPROXY:"${GOPROXY}

moduleName=$(go list -m)
program=$(basename ${moduleName})

cd ${projectpath} || exit 1

go mod tidy || exit 1

# sqlite 不支持交叉编译
# Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work
go build -v -ldflags "$flags" -o ${projectpath}/$program  || exit 1

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
cp ${projectpath}/config.yml ${BUILD_TARGET_DIR}/ || exit 1

# 打包

# make it execute able
chmod +x ${BUILD_TARGET_DIR}/server

os_name=$(uname -s)
typeset -l os_name

# tar the jrasp-server.tar.gz
cd ${BUILD_TMP} || exit_on_err 1 "[ERROR] cd target dir error."

zip_name="jrasp-server-bin-$(date +"%Y%m%d%H%M%S").tar.gz"

tar -zcvf ${zip_name} server/

echo "$zip_name finish."

exit 0
