package report

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"github.com/gookit/goutil/fsutil"
	"github.com/imroc/req"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"server/common"
	"server/config"
	"server/model"
	"server/util"
)

func (this *UpdateClient) handleClientUpgrade(messageContent ClientUpgradeRequest) {
	// 参数校验
	if err := common.Validate.Struct(&messageContent); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		if err != nil {
			common.Log.Error("校验数据格式出错 err: %s", errStr)
		}
		return
	}
	switch messageContent.Type {
	case "config":
		this.upgradeConfig(messageContent)
	case "server":
		success := this.upgradeServer(messageContent)
		if success {
			os.Exit(0) // 进程退出
		}
	default:
	}
}

func (this *UpdateClient) upgradeConfig(messageContent ClientUpgradeRequest) {
	defer func() {
		upgradeResult := WebSocketMessageRequest{
			MessageType:    CLIENT_UPGRADE_RESULT,
			MessageContent: structs.Map(messageContent),
		}
		err := this.conn.WriteJSON(upgradeResult)
		if err != nil {
			common.Log.Errorf("发送心跳数据失败, error: %v", err)
		}
	}()

	resp, err := req.Get(messageContent.DownloadUrl)
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
		return
	}
	downPath := filepath.Join(config.Conf.Env.WorkDir, "tmp")
	if !fsutil.PathExists(downPath) {
		err = os.MkdirAll(downPath, 0755)
		if err != nil {
			messageContent.State = 2
			messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
			return
		}
	}
	newFilePath := filepath.Join(downPath, "config.zip")
	err = resp.ToFile(newFilePath)
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
		return
	}
	newMd5, err := util.GetFileMd5(newFilePath)
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
		return
	}
	if newMd5 != messageContent.Md5 {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", "校验md5失败")
		return
	}

	// 开始更新通用防护策略
	r, err := zip.OpenReader(newFilePath)
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
		return
	}
	// 读取config.json配置文件并反序列化
	configHandler, err := r.Open("config.json")
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", "读取config.json失败")
		return
	}
	jsonData, err := io.ReadAll(configHandler)
	configHandler.Close()
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", "读取config.json失败")
		return
	}
	var exportConfig model.RaspExportConfig
	err = json.Unmarshal(jsonData, &exportConfig)
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", "反序列化配置失败")
		return
	}
	// 写入file信息
	for _, item := range exportConfig.RaspFileInfo {
		fileInfo, err := this.RaspFileRepository.GetRaspFileByHash(item.FileHash)
		if err != nil {
			messageContent.State = 2
			messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
			return
		}
		// 如果文件不存在则写入
		if fileInfo == nil {
			fileHandler, err := r.Open(path.Join("file", item.FileName))
			if err != nil {
				messageContent.State = 2
				messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
				return
			}
			data, _ := io.ReadAll(fileHandler)
			actualHash, err := util.GetMd5FromBytes(data)
			if actualHash != item.FileHash {
				messageContent.State = 2
				messageContent.Message = fmt.Sprintf("校验文件hash失败, actual hash:%v, expect hash: %v", actualHash, item.FileHash)
				return
			}
			// hash校验成功
			uploadDir := path.Dir(item.DiskPath)
			if !fsutil.PathExist(uploadDir) {
				err = os.MkdirAll(uploadDir, 0755)
				if err != nil {
					messageContent.State = 2
					messageContent.Message = fmt.Sprintf("创建上传目录失败, %v", err)
					return
				}
			}
			fileHandler.Close()
			err = os.WriteFile(item.DiskPath, data, 0755)
			if err != nil {
				messageContent.State = 2
				messageContent.Message = fmt.Sprintf("写入文件记录失败, %v", err)
				return
			}
			// 构建表file记录
			raspFile := &model.RaspFile{
				Timestamp:   item.Timestamp,
				FileName:    item.FileName,
				FileHash:    item.FileHash,
				DiskPath:    item.DiskPath,
				DownLoadUrl: item.DownLoadUrl,
				MimeType:    item.MimeType,
				Creator:     "system",
			}
			err = this.RaspFileRepository.CreateRaspFile(raspFile)
			if err != nil {
				messageContent.State = 2
				messageContent.Message = fmt.Sprintf("创建文件记录失败, %v", err)
				return
			}
		} else {
			common.Log.Debugf("文件: %v 已存在, 跳过导入", item.FileName)
		}
	}
	// 写入module信息
	for _, item := range exportConfig.RaspModuleInfo {
		moduleInfo, err := this.RaspModuleRepository.GetRaspModuleByName(item.ModuleName)
		if err != nil {
			messageContent.State = 2
			messageContent.Message = fmt.Sprintf("查询module信息失败, %v", err)
			return
		}
		// 如果模块名不存在则直接创建
		if moduleInfo == nil {
			raspModule := &model.RaspModule{
				RowGuid:       item.RowGuid,
				ModuleName:    item.ModuleName,
				ModuleVersion: item.ModuleVersion,
				Desc:          item.Desc,
				Parameters:    item.Parameters,
				Upgradable:    item.Upgradable,
				Creator:       "system",
				Operator:      "system",
			}
			err = this.RaspModuleRepository.CreateRaspModule(raspModule)
			if err != nil {
				messageContent.State = 2
				messageContent.Message = fmt.Sprintf("新建module记录失败, %v", err)
				return
			}
		} else {
			// 如果模块名存在则替换版本、描述、默认用户参数
			moduleInfo.RowGuid = item.RowGuid
			moduleInfo.ModuleName = item.ModuleName
			moduleInfo.ModuleVersion = item.ModuleVersion
			moduleInfo.Desc = item.Desc
			moduleInfo.Parameters = item.Parameters
			moduleInfo.Upgradable = item.Upgradable
			moduleInfo.Creator = "system"
			moduleInfo.Operator = "system"
			err = this.RaspModuleRepository.UpdateRaspModule(moduleInfo)
			if err != nil {
				messageContent.State = 2
				messageContent.Message = fmt.Sprintf("更新module信息失败, %v", err)
				return
			}
		}
	}
	// 写入component信息
	for _, item := range exportConfig.RaspComponentInfo {
		componentInfo, err := this.RaspComponentRepository.GetRaspComponentsByGuidAndName(item.ParentGuid, item.ComponentName)
		if err != nil {
			messageContent.State = 2
			messageContent.Message = fmt.Sprintf("查询component信息失败, %v", err)
			return
		}
		if componentInfo == nil {
			raspComponent := &model.RaspComponent{
				ParentGuid:       item.ParentGuid,
				ComponentName:    item.ComponentName,
				ComponentType:    item.ComponentType,
				ComponentVersion: item.ComponentVersion,
				DownLoadURL:      item.DownLoadURL,
				Md5:              item.Md5,
				Parameters:       item.Parameters,
			}
			err = this.RaspComponentRepository.CreateRaspComponent(raspComponent)
			if err != nil {
				messageContent.State = 2
				messageContent.Message = fmt.Sprintf("新建component记录失败, %v", err)
				return
			}
		} else {
			componentInfo.ParentGuid = item.ParentGuid
			componentInfo.ComponentName = item.ComponentName
			componentInfo.ComponentType = item.ComponentType
			componentInfo.ComponentVersion = item.ComponentVersion
			componentInfo.DownLoadURL = item.DownLoadURL
			componentInfo.Md5 = item.Md5
			componentInfo.Parameters = item.Parameters
			err = this.RaspComponentRepository.UpdateRaspComponent(componentInfo)
			if err != nil {
				messageContent.State = 2
				messageContent.Message = fmt.Sprintf("更新module信息失败, %v", err)
				return
			}
		}
	}
	// 写入config信息
	configInfo, err := this.RaspConfigRepository.GetRaspConfigByName(exportConfig.RaspConfigInfo.Name)
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("查询config信息失败, %v", err)
		return
	}
	if configInfo == nil {
		raspConfig := &model.RaspConfig{
			RowGuid:   exportConfig.RaspConfigInfo.RowGuid,
			Name:      exportConfig.RaspConfigInfo.Name,
			Version:   1,
			Desc:      exportConfig.RaspConfigInfo.Desc,
			Status:    exportConfig.RaspConfigInfo.Status,
			Creator:   "system",
			Operator:  "system",
			IsDefault: exportConfig.RaspConfigInfo.IsDefault,
		}
		err = this.RaspConfigRepository.CreateRaspConfig(raspConfig)
		if err != nil {
			messageContent.State = 2
			messageContent.Message = fmt.Sprintf("创建config信息失败, %v", err)
			return
		}
		// 写入历史版本信息
		raspConfigHistory := &model.RaspConfigHistory{
			ParentGuid:    raspConfig.RowGuid,
			Version:       1,
			AgentMode:     exportConfig.RaspConfigHistoryInfo.AgentMode,
			ModuleConfigs: exportConfig.RaspConfigHistoryInfo.ModuleConfigs,
			LogPath:       exportConfig.RaspConfigHistoryInfo.LogPath,
			AgentConfigs:  exportConfig.RaspConfigHistoryInfo.AgentConfigs,
			RaspBinInfo:   exportConfig.RaspConfigHistoryInfo.RaspBinInfo,
			RaspLibInfo:   exportConfig.RaspConfigHistoryInfo.RaspLibInfo,
			Desc:          exportConfig.RaspConfigHistoryInfo.Desc,
			Source:        "导入版本",
		}
		err = this.RaspConfigHistoryRepository.CreateRaspConfigHistory(raspConfigHistory)
		if err != nil {
			messageContent.State = 2
			messageContent.Message = fmt.Sprintf("创建config版本信息失败, %v", err)
			return
		}
	} else {
		configInfo.Name = exportConfig.RaspConfigInfo.Name
		configInfo.Version += 1
		configInfo.Desc = exportConfig.RaspConfigInfo.Desc
		configInfo.Status = exportConfig.RaspConfigInfo.Status
		configInfo.Creator = "system"
		configInfo.Operator = "system"
		configInfo.IsDefault = exportConfig.RaspConfigInfo.IsDefault
		err = this.RaspConfigRepository.UpdateRaspConfig(configInfo)
		if err != nil {
			messageContent.State = 2
			messageContent.Message = fmt.Sprintf("更新config信息失败, %v", err)
			return
		}
		// 写入历史版本信息
		raspConfigHistory := &model.RaspConfigHistory{
			ParentGuid:    configInfo.RowGuid,
			Version:       configInfo.Version,
			AgentMode:     exportConfig.RaspConfigHistoryInfo.AgentMode,
			ModuleConfigs: exportConfig.RaspConfigHistoryInfo.ModuleConfigs,
			LogPath:       exportConfig.RaspConfigHistoryInfo.LogPath,
			AgentConfigs:  exportConfig.RaspConfigHistoryInfo.AgentConfigs,
			RaspBinInfo:   exportConfig.RaspConfigHistoryInfo.RaspBinInfo,
			RaspLibInfo:   exportConfig.RaspConfigHistoryInfo.RaspLibInfo,
			Desc:          exportConfig.RaspConfigHistoryInfo.Desc,
			Source:        "导入版本",
		}
		err = this.RaspConfigHistoryRepository.CreateRaspConfigHistory(raspConfigHistory)
		if err != nil {
			messageContent.State = 2
			messageContent.Message = fmt.Sprintf("创建config版本信息失败, %v", err)
			return
		}
	}
	// 更新成功
	messageContent.State = 1
	messageContent.Message = "更新成功"
	return
}

func (this *UpdateClient) upgradeServer(messageContent ClientUpgradeRequest) bool {
	defer func() {
		upgradeResult := WebSocketMessageRequest{
			MessageType:    CLIENT_UPGRADE_RESULT,
			MessageContent: structs.Map(messageContent),
		}
		err := this.conn.WriteJSON(upgradeResult)
		if err != nil {
			common.Log.Errorf("发送心跳数据失败, error: %v", err)
		}
	}()

	resp, err := req.Get(messageContent.DownloadUrl)
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
		return false
	}
	downPath := filepath.Join(config.Conf.Env.WorkDir, "tmp")
	if !fsutil.PathExists(downPath) {
		err = os.MkdirAll(downPath, 0755)
		if err != nil {
			messageContent.State = 2
			messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
			return false
		}
	}
	newFilePath := filepath.Join(downPath, "server.zip")
	err = resp.ToFile(newFilePath)
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
		return false
	}
	newMd5, err := util.GetFileMd5(newFilePath)
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
		return false
	}
	if newMd5 != messageContent.Md5 {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", "校验md5失败")
		return false
	}
	// 开始更新
	data, err := util.ReadFileFromZipByPath(newFilePath, path.Join(runtime.GOOS, runtime.GOARCH, config.Conf.Env.BinFileName))
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
		return false
	}
	// 兼容windows环境下直接写daemon导致文件被占用的错误，应该先把运行中的daemon重命名为临时文件，再将新的文件覆盖
	err = os.Rename(config.Conf.Env.BinFileName, config.Conf.Env.BinFileName+".del")
	if err != nil {
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
		return false
	}
	err = os.WriteFile(config.Conf.Env.BinFileName, data, 0777)
	if err != nil {
		common.Log.Errorf("replace server error: %v", err)
		messageContent.State = 2
		messageContent.Message = fmt.Sprintf("更新失败, error: %v", err)
		return false
	} else {
		common.Log.Infof("update server file success, process will exit...")
		this.clean(newFilePath)
		exist := fsutil.PathExists(config.Conf.Env.BinFileName)
		if exist {
			// 更新成功
			messageContent.State = 1
			messageContent.Message = "更新成功"
			return true
		} else {
			common.Log.Error("server 文件不存在")
			messageContent.State = 2
			messageContent.Message = "server 文件不存在"
			return false
		}
	}
}

func (this *UpdateClient) clean(filePath string) {
	exists := fsutil.PathExists(filePath)
	if exists {
		err := os.Remove(filePath)
		if err != nil {
			common.Log.Errorf("delete file error: %v", err)
		}
	}
}
