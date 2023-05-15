package controller

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gookit/goutil/fsutil"
	"github.com/mholt/archiver/v4"
	"io"
	"os"
	"path"
	"path/filepath"
	"server/common"
	"server/model"
	"server/repository"
	"server/response"
	"server/util"
	"server/vo"
	"strings"
)

type IRaspConfigController interface {
	CreateRaspConfig(c *gin.Context)
	UpdateRaspConfig(c *gin.Context)
	GetRaspConfigs(c *gin.Context)
	BatchDeleteConfigByIds(c *gin.Context)
	GetViperRaspConfig(c *gin.Context)
	UpdateRaspConfigStatusById(c *gin.Context)
	UpdateRaspConfigDefaultById(c *gin.Context)
	PushRaspConfig(c *gin.Context)
	CopyRaspConfig(c *gin.Context)
	GetRaspModules(c *gin.Context)
	ExportRaspConfig(c *gin.Context)
	ImportRaspConfig(c *gin.Context)
}

type RaspConfigController struct {
	RaspConfigRepository repository.IRaspConfigRepository
	RaspModuleRepository repository.IRaspModuleRepository
	RaspFileRepository   repository.IRaspFileRepository
}

func NewRaspConfigController() IRaspConfigController {
	raspConfigRepository := repository.NewRaspConfigRepository()
	raspModuleRepository := repository.NewRaspModuleRepository()
	raspFileRepository := repository.NewRaspFileRepository()
	raspConfigController := RaspConfigController{
		RaspConfigRepository: raspConfigRepository,
		RaspModuleRepository: raspModuleRepository,
		RaspFileRepository:   raspFileRepository,
	}
	return raspConfigController
}

func (r RaspConfigController) GetRaspConfigs(c *gin.Context) {
	var req vo.RaspConfigListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	// 获取
	raspConfigs, total, err := r.RaspConfigRepository.GetRaspConfigs(&req)
	if err != nil {
		response.Fail(c, nil, "获取配置列表失败")
		return
	}
	response.Success(c, gin.H{
		"list": raspConfigs, "total": total,
	}, "获取配置列表成功")
}

func (r RaspConfigController) CreateRaspConfig(c *gin.Context) {
	var req vo.CreateRaspConfigRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// 获取当前用户
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}

	raspConfig := model.RaspConfig{
		Name:          req.Name,
		Desc:          req.Desc,
		Status:        req.Status,
		Creator:       ctxUser.Username,
		Operator:      ctxUser.Username,
		AgentMode:     req.AgentMode,
		ModuleConfigs: req.ModuleConfigs,
		LogPath:       req.LogPath,
		AgentConfigs:  req.AgentConfigs,
		RaspBinInfo:   req.RaspBinInfo,
		RaspLibInfo:   req.RaspLibInfo,
	}

	// 获取
	err = r.RaspConfigRepository.CreateRaspConfig(&raspConfig)
	if err != nil {
		response.Fail(c, nil, "创建配置失败"+err.Error())
		return
	}
	response.Success(c, nil, "创建配置成功")
	return
}

func (r RaspConfigController) UpdateRaspConfig(c *gin.Context) {
	var req vo.UpdateRaspConfigRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// 获取当前用户
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}

	id := req.ID
	config, err := r.RaspConfigRepository.GetRaspConfigById(id)
	if err != nil {
		response.Fail(c, nil, "获取当前模块失败")
		return
	}

	config.Name = req.Name
	config.Desc = req.Desc
	config.Status = req.Status
	config.Operator = ctxUser.Username
	config.AgentMode = req.AgentMode
	config.ModuleConfigs = req.ModuleConfigs
	config.LogPath = req.LogPath
	config.AgentConfigs = req.AgentConfigs
	config.RaspBinInfo = req.RaspBinInfo
	config.RaspLibInfo = req.RaspLibInfo

	err = r.RaspConfigRepository.UpdateRaspConfig(config)
	if err != nil {
		response.Fail(c, nil, "更新当前配置失败")
		return
	}
	response.Success(c, nil, "更新配置成功")
}

// 批量删除接口
func (r RaspConfigController) BatchDeleteConfigByIds(c *gin.Context) {
	var req vo.DeleteRaspConfigRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// 删除接口
	err := r.RaspConfigRepository.DeleteRaspConfig(req.Ids)
	if err != nil {
		response.Fail(c, nil, "删除配置失败: "+err.Error())
		return
	}
	response.Success(c, nil, "删除配置成功")
}

// GetViperRaspConfig viper remote get
func (l RaspConfigController) GetViperRaspConfig(c *gin.Context) {
	var req vo.RaspConfigRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	name := strings.TrimSpace(req.Key)
	config, err := l.RaspConfigRepository.GetRaspConfigByName(name)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"key": name, "value": config.AgentConfigs}, "获取配置成功")
}

func (l RaspConfigController) UpdateRaspConfigStatusById(c *gin.Context) {
	var req vo.UpdateRaspConfigStatusRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	// 获取当前用户
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}

	id := req.ID
	config, err := l.RaspConfigRepository.GetRaspConfigById(id)
	if err != nil {
		response.Fail(c, nil, "更新当前策略失败")
		return
	}
	config.Status = !config.Status
	config.Operator = ctxUser.Username

	err = l.RaspConfigRepository.UpdateRaspConfig(config)
	if err != nil {
		response.Fail(c, nil, "更新当前策略失败")
		return
	}
	response.Success(c, nil, "更新策略成功")
}

func (l RaspConfigController) UpdateRaspConfigDefaultById(c *gin.Context) {
	var req vo.UpdateRaspConfigDefaultRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	// 全部设置为非默认
	if err := common.DB.Model(&model.RaspConfig{}).Where("1=1").Updates(map[string]interface{}{"is_default": false}).Error; err != nil {
		response.Fail(c, nil, "更新默认策略失败")
		return
	}

	// 获取当前用户
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}
	id := req.ID
	config, err := l.RaspConfigRepository.GetRaspConfigById(id)
	if err != nil {
		response.Fail(c, nil, "更新当前策略失败")
		return
	}
	config.IsDefault = req.IsDefault
	config.Operator = ctxUser.Username

	err = l.RaspConfigRepository.UpdateRaspConfig(config)
	if err != nil {
		response.Fail(c, nil, "更新当前策略失败")
		return
	}
	response.Success(c, nil, "更新策略成功")
}

func (l RaspConfigController) PushRaspConfig(c *gin.Context) {
	var req vo.PushRaspConfigRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	// 生成推送的配置信息
	hostController := NewRaspHostController()
	hostRepository := repository.NewRaspHostRepository()
	content, err := hostController.GeneratePushConfig(req.ID)
	if err != nil {
		response.Fail(c, nil, "推送配置失败")
		return
	}
	// 先更新数据库中的配置Id
	hostList, _, err := hostRepository.GetRaspHosts(new(vo.RaspHostListRequest))
	var hostNameList []string
	for _, item := range hostList {
		item.ConfigId = req.ID
		err = hostRepository.UpdateRaspHost(item)
		if err != nil {
			common.Log.Errorf("更新配置Id失败, %v", err)
			continue
		}
		hostNameList = append(hostNameList, item.HostName)
	}
	// 开始推送最新的配置
	hostController.PushHostsConfig(hostNameList, content)

	response.Success(c, nil, "推送策略成功")
}

func (l RaspConfigController) CopyRaspConfig(c *gin.Context) {
	var req vo.CopyRaspConfigRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	// 获取要复制的配置对象
	srcConfig, err := l.RaspConfigRepository.GetRaspConfigById(req.ID)
	if err != nil {
		response.Fail(c, nil, "复制当前配置失败")
		return
	}
	// 获取当前用户
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}
	destConfig := model.RaspConfig{
		Name:          srcConfig.Name + "_Copy",
		Desc:          srcConfig.Desc,
		Status:        srcConfig.Status,
		Creator:       ctxUser.Username,
		Operator:      ctxUser.Username,
		AgentMode:     srcConfig.AgentMode,
		ModuleConfigs: srcConfig.ModuleConfigs,
		LogPath:       srcConfig.LogPath,
		AgentConfigs:  srcConfig.AgentConfigs,
		RaspBinInfo:   srcConfig.RaspBinInfo,
		RaspLibInfo:   srcConfig.RaspLibInfo,
		IsDefault:     false,
	}
	err = l.RaspConfigRepository.CreateRaspConfig(&destConfig)
	if err != nil {
		response.Fail(c, nil, "复制当前配置失败"+err.Error())
		return
	}
	response.Success(c, nil, "复制配置成功")
	return
}

func (l RaspConfigController) GetRaspModules(c *gin.Context) {
	var req vo.RaspCheckboxModuleListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	modules, total, err := l.RaspConfigRepository.GetRaspModules()
	if err != nil {
		response.Fail(c, nil, "复制当前配置失败"+err.Error())
		return
	}
	var results []map[string]interface{}
	for _, module := range modules {
		var item = make(map[string]interface{})
		item["moduleName"] = module.ModuleName
		item["moduleVersion"] = module.ModuleVersion
		results = append(results, item)
	}
	response.Success(c, gin.H{
		"list": results, "total": total,
	}, "复制配置成功")
}

func (l RaspConfigController) ExportRaspConfig(c *gin.Context) {
	var req vo.ExportRaspConfigRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	// 创建导出临时文件夹
	exportDir := filepath.Join("export")
	if !fsutil.PathExist(exportDir) {
		err := os.Mkdir(exportDir, 0755)
		if err != nil {
			response.Fail(c, nil, "导出当前配置失败")
			return
		}
	}
	// 获取要导出的配置对象
	srcConfig, err := l.RaspConfigRepository.GetRaspConfigById(req.ID)
	if err != nil {
		response.Fail(c, nil, "获取当前配置失败")
		return
	}
	// 获取module信息
	var moduleConfigs []model.ModuleConfig
	var libConfigs map[string]interface{}
	var binConfigs map[string]interface{}
	var raspModules []model.RaspModule
	err = json.Unmarshal(srcConfig.ModuleConfigs, &moduleConfigs)
	if err != nil {
		response.Fail(c, nil, "导出当前配置失败")
		return
	}
	err = json.Unmarshal(srcConfig.RaspBinInfo, &binConfigs)
	if err != nil {
		response.Fail(c, nil, "导出当前配置失败")
		return
	}
	err = json.Unmarshal(srcConfig.RaspLibInfo, &libConfigs)
	if err != nil {
		response.Fail(c, nil, "导出当前配置失败")
		return
	}
	for _, item := range moduleConfigs {
		moduleInfo, err := l.RaspModuleRepository.GetRaspModuleByName(item.ModuleName, item.ModuleVersion)
		if err != nil {
			response.Fail(c, nil, "导出当前配置失败")
			return
		}
		raspModules = append(raspModules, *moduleInfo)
	}
	// 获取file信息
	var raspFiles []model.RaspFile
	for _, item := range raspModules {
		fileInfo, err := l.RaspFileRepository.GetRaspFileByHash(item.Md5)
		if err != nil {
			response.Fail(c, nil, "导出当前配置失败")
			return
		}
		raspFiles = append(raspFiles, *fileInfo)
	}
	fileInfo, err := l.RaspFileRepository.GetRaspFileByHash(libConfigs["md5"].(string))
	if err != nil {
		response.Fail(c, nil, "导出当前配置失败")
		return
	}
	raspFiles = append(raspFiles, *fileInfo)
	fileInfo, err = l.RaspFileRepository.GetRaspFileByHash(binConfigs["md5"].(string))
	if err != nil {
		response.Fail(c, nil, "导出当前配置失败")
		return
	}
	raspFiles = append(raspFiles, *fileInfo)
	// 构件配置对象
	var exportConfig model.RaspExportConfig
	exportConfig.RaspConfigInfo = *srcConfig
	exportConfig.RaspModuleInfo = raspModules
	exportConfig.RaspFileInfo = raspFiles

	jsonStr, err := json.Marshal(exportConfig)
	if err != nil {
		response.Fail(c, nil, "序列化当前配置失败")
		return
	}
	err = os.WriteFile(filepath.Join(exportDir, "config.json"), jsonStr, 0755)
	if err != nil {
		response.Fail(c, nil, "写当前配置到服务器失败")
		return
	}
	// 导出jar包和zip文件
	var fileMaps = make(map[string]string)
	if !fsutil.PathExist(filepath.Join(exportDir, "file")) {
		err = os.Mkdir(filepath.Join(exportDir, "file"), 0777)
	}
	for _, item := range exportConfig.RaspFileInfo {
		exportPath := filepath.Join(exportDir, "file", item.FileName)
		err = fsutil.CopyFile(item.DiskPath, exportPath)
		if err != nil {
			response.Fail(c, nil, "写当前配置到服务器失败")
			return
		}
		fileMaps[exportPath] = path.Join("file", item.FileName)
	}
	fileMaps[filepath.Join("export", "config.json")] = "config.json"
	files, err := archiver.FilesFromDisk(nil, fileMaps)
	if err != nil {
		response.Fail(c, nil, "压缩配置到服务器失败")
		return
	}

	format := archiver.Zip{}
	var data bytes.Buffer
	err = format.Archive(context.Background(), &data, files)
	if err != nil {
		response.Fail(c, nil, "压缩配置到服务器失败")
		return
	}
	err = os.RemoveAll(exportDir)
	if err != nil {
		response.Fail(c, nil, "删除临时失败")
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fmt.Sprintf("export_%v.zip", srcConfig.Name)))
	c.Writer.Header().Add("response-type", "blob")
	c.Data(200, "application/octet-stream", data.Bytes())
}

func (l RaspConfigController) ImportRaspConfig(c *gin.Context) {
	// 获取当前用户
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["files"]
	file := files[0]
	zipHandler, err := file.Open()
	if err != nil {
		response.Fail(c, nil, "导入配置失败")
		return
	}
	defer zipHandler.Close()
	r, err := zip.NewReader(zipHandler, file.Size)
	if err != nil {
		response.Fail(c, nil, "导入配置失败")
		return
	}
	// 读取config.json配置文件并反序列化
	configHandler, err := r.Open("config.json")
	if err != nil {
		response.Fail(c, nil, "读取config.json失败")
		return
	}
	jsonData, err := io.ReadAll(configHandler)
	configHandler.Close()
	if err != nil {
		response.Fail(c, nil, "读取config.json失败")
		return
	}
	var exportConfig model.RaspExportConfig
	err = json.Unmarshal(jsonData, &exportConfig)
	if err != nil {
		response.Fail(c, nil, "反序列化配置失败")
		return
	}
	// 写入file信息
	for _, item := range exportConfig.RaspFileInfo {
		fileInfo, err := l.RaspFileRepository.GetRaspFileByHash(item.FileHash)
		if err != nil {
			response.Fail(c, nil, fmt.Sprintf("查询file信息失败, %v", err))
			return
		}
		// 如果文件不存在则写入
		if fileInfo == nil {
			fileHandler, err := r.Open(path.Join("file", item.FileName))
			if err != nil {
				response.Fail(c, nil, fmt.Sprintf("读取文件失败, %v", err))
				return
			}
			data, _ := io.ReadAll(fileHandler)
			actualHash, err := util.GetMd5FromBytes(data)
			if actualHash != item.FileHash {
				response.Fail(c, nil, fmt.Sprintf("校验文件hash失败, actual hash:%v, expect hash: %v", actualHash, item.FileHash))
				return
			}
			// hash校验成功
			uploadDir := path.Dir(item.DiskPath)
			if !fsutil.PathExist(uploadDir) {
				err = os.MkdirAll(uploadDir, 0755)
				if err != nil {
					response.Fail(c, nil, fmt.Sprintf("创建上传目录失败, %v", err))
					return
				}
			}
			fileHandler.Close()
			err = os.WriteFile(item.DiskPath, data, 0755)
			if err != nil {
				response.Fail(c, nil, fmt.Sprintf("写入文件记录失败, %v", err))
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
				Creator:     ctxUser.Username,
			}
			err = l.RaspFileRepository.CreateRaspFile(raspFile)
			if err != nil {
				response.Fail(c, nil, fmt.Sprintf("创建文件记录失败, %v", err))
				return
			}
		} else {
			common.Log.Debugf("文件: %v 已存在, 跳过导入", item.FileName)
		}
	}
	// 写入module信息
	for _, item := range exportConfig.RaspModuleInfo {
		moduleInfo, err := l.RaspModuleRepository.GetRaspModuleByName(item.ModuleName, item.ModuleVersion)
		if err != nil {
			response.Fail(c, nil, fmt.Sprintf("查询module信息失败, %v", err))
			return
		}
		if moduleInfo == nil {
			raspModule := &model.RaspModule{
				ModuleName:    item.ModuleName,
				ModuleType:    item.ModuleType,
				ModuleVersion: item.ModuleVersion,
				DownLoadURL:   item.DownLoadURL,
				Md5:           item.Md5,
				Desc:          item.Desc,
				Status:        item.Status,
				Parameters:    item.Parameters,
				Creator:       ctxUser.Username,
				Operator:      ctxUser.Username,
			}
			err = l.RaspModuleRepository.CreateRaspModule(raspModule)
			if err != nil {
				response.Fail(c, nil, fmt.Sprintf("新建module记录失败, %v", err))
				return
			}
		} else {
			moduleInfo.ModuleName = item.ModuleName
			moduleInfo.ModuleType = item.ModuleType
			moduleInfo.ModuleVersion = item.ModuleVersion
			moduleInfo.DownLoadURL = item.DownLoadURL
			moduleInfo.Md5 = item.Md5
			moduleInfo.Desc = item.Desc
			moduleInfo.Status = item.Status
			moduleInfo.Parameters = item.Parameters
			moduleInfo.Creator = ctxUser.Username
			moduleInfo.Operator = ctxUser.Username
			err = l.RaspModuleRepository.UpdateRaspModule(moduleInfo)
			if err != nil {
				response.Fail(c, nil, fmt.Sprintf("更新module信息失败, %v", err))
				return
			}
		}
	}
	// 写入config信息
	configInfo, err := l.RaspConfigRepository.GetRaspConfigByName(exportConfig.RaspConfigInfo.Name)
	if err != nil {
		response.Fail(c, nil, fmt.Sprintf("查询config信息失败, %v", err))
		return
	}
	if configInfo == nil {
		raspConfig := &model.RaspConfig{
			Name:          exportConfig.RaspConfigInfo.Name,
			Desc:          exportConfig.RaspConfigInfo.Desc,
			Status:        exportConfig.RaspConfigInfo.Status,
			Creator:       ctxUser.Username,
			Operator:      ctxUser.Username,
			AgentMode:     exportConfig.RaspConfigInfo.AgentMode,
			ModuleConfigs: exportConfig.RaspConfigInfo.ModuleConfigs,
			LogPath:       exportConfig.RaspConfigInfo.LogPath,
			AgentConfigs:  exportConfig.RaspConfigInfo.AgentConfigs,
			RaspBinInfo:   exportConfig.RaspConfigInfo.RaspBinInfo,
			RaspLibInfo:   exportConfig.RaspConfigInfo.RaspLibInfo,
			IsDefault:     exportConfig.RaspConfigInfo.IsDefault,
		}
		err = l.RaspConfigRepository.CreateRaspConfig(raspConfig)
		if err != nil {
			response.Fail(c, nil, fmt.Sprintf("创建config信息失败, %v", err))
			return
		}
	} else {
		configInfo.Name = exportConfig.RaspConfigInfo.Name
		configInfo.Desc = exportConfig.RaspConfigInfo.Desc
		configInfo.Status = exportConfig.RaspConfigInfo.Status
		configInfo.Creator = ctxUser.Username
		configInfo.Operator = ctxUser.Username
		configInfo.AgentMode = exportConfig.RaspConfigInfo.AgentMode
		configInfo.ModuleConfigs = exportConfig.RaspConfigInfo.ModuleConfigs
		configInfo.LogPath = exportConfig.RaspConfigInfo.LogPath
		configInfo.AgentConfigs = exportConfig.RaspConfigInfo.AgentConfigs
		configInfo.RaspBinInfo = exportConfig.RaspConfigInfo.RaspBinInfo
		configInfo.RaspLibInfo = exportConfig.RaspConfigInfo.RaspLibInfo
		configInfo.IsDefault = exportConfig.RaspConfigInfo.IsDefault
		err = l.RaspConfigRepository.UpdateRaspConfig(configInfo)
		if err != nil {
			response.Fail(c, nil, fmt.Sprintf("更新config信息失败, %v", err))
			return
		}
	}
	response.Success(c, nil, "导入配置成功")
}
