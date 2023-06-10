package common

import (
	"os"
	"path/filepath"
	"server/config"
	"strconv"
	"strings"
)

func InitInstallScripts() {
	data, err := os.ReadFile(filepath.Join("install", "install-agent-template.sh"))
	if err != nil {
		Log.Errorf("读取文件: %v 失败", "install-agent-template.sh")
		return
	}
	text := string(data)
	text = strings.Replace(text, "{RASP_SERVER_IP}", config.Conf.Env.Ip, 1)
	text = strings.Replace(text, "{RASP_SERVER_PORT}", strconv.Itoa(config.Conf.System.Port), 1)
	err = os.WriteFile(filepath.Join("install", "install-agent.sh"), []byte(text), 0777)
	if err != nil {
		Log.Errorf("写入文件: %v 失败", "install-agent.sh")
		return
	}
	Log.Infof("初始化安装脚本成功")
}
