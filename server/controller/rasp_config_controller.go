package controller

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gookit/goutil/fsutil"
	"github.com/mholt/archiver/v4"
	"gorm.io/datatypes"
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
	SyncRaspConfig(c *gin.Context)
}

type RaspConfigController struct {
	RaspConfigRepository        repository.IRaspConfigRepository
	RaspConfigHistoryRepository repository.IRaspConfigHistoryRepository
	RaspModuleRepository        repository.IRaspModuleRepository
	RaspComponentRepository     repository.IRaspComponentRepository
	RaspFileRepository          repository.IRaspFileRepository
}

func NewRaspConfigController() IRaspConfigController {
	repo1 := repository.NewRaspConfigRepository()
	repo2 := repository.NewRaspModuleRepository()
	repo3 := repository.NewRaspFileRepository()
	repo4 := repository.NewRaspComponentRepository()
	repo5 := repository.NewRaspConfigHistoryRepository()
	raspConfigController := RaspConfigController{
		RaspConfigRepository:        repo1,
		RaspModuleRepository:        repo2,
		RaspFileRepository:          repo3,
		RaspComponentRepository:     repo4,
		RaspConfigHistoryRepository: repo5,
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
		RowGuid:  uuid.New().String(),
		Name:     req.Name,
		Version:  req.Version + 1,
		Desc:     req.Desc,
		Status:   req.Status,
		Creator:  ctxUser.Username,
		Operator: ctxUser.Username,
	}

	newModuleConfigs, err := r.generateModuleConfig(req.ModuleConfigs)
	if err != nil {
		response.Fail(c, nil, "获取模块最新下载地址出错")
		return
	}
	raspConfigHistory := model.RaspConfigHistory{
		ParentGuid:    raspConfig.RowGuid,
		Version:       req.Version + 1,
		AgentMode:     req.AgentMode,
		ModuleConfigs: newModuleConfigs,
		LogPath:       req.LogPath,
		AgentConfigs:  req.AgentConfigs,
		RaspBinInfo:   req.RaspBinInfo,
		RaspLibInfo:   req.RaspLibInfo,
		Desc:          req.HistoryDesc,
		Source:        "新建版本",
	}

	// 获取
	err = r.RaspConfigRepository.CreateRaspConfig(&raspConfig)
	if err != nil {
		response.Fail(c, nil, "创建配置失败"+err.Error())
		return
	}
	err = r.RaspConfigHistoryRepository.CreateRaspConfigHistory(&raspConfigHistory)
	if err != nil {
		response.Fail(c, nil, "创建配置版本失败"+err.Error())
		return
	}
	response.Success(c, nil, "创建配置成功")
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
	// 默认防护策略仅允许root用户更新
	if id == 1 && ctxUser.Username != "root" {
		response.Fail(c, nil, "通用防护策略不允许更新，请自行新建策略")
		return
	}

	config.Name = req.Name
	config.Version = util.Ternary(req.IsNewVersion, req.Version+1, req.Version).(int)
	config.Desc = req.Desc
	config.Status = req.Status
	config.Operator = ctxUser.Username
	config.IsModified = true

	err = r.RaspConfigRepository.UpdateRaspConfig(config)
	if err != nil {
		response.Fail(c, nil, "更新当前配置失败")
		return
	}
	if req.IsNewVersion {
		newModuleConfigs, err := r.generateModuleConfig(req.ModuleConfigs)
		if err != nil {
			response.Fail(c, nil, "获取模块最新下载地址出错")
			return
		}
		raspConfigHistory := model.RaspConfigHistory{
			ParentGuid:    config.RowGuid,
			Version:       req.Version + 1,
			AgentMode:     req.AgentMode,
			ModuleConfigs: newModuleConfigs,
			LogPath:       req.LogPath,
			AgentConfigs:  req.AgentConfigs,
			RaspBinInfo:   req.RaspBinInfo,
			RaspLibInfo:   req.RaspLibInfo,
			Desc:          req.HistoryDesc,
			Source:        "更新版本",
		}
		err = r.RaspConfigHistoryRepository.CreateRaspConfigHistory(&raspConfigHistory)
		if err != nil {
			response.Fail(c, nil, "创建配置版本失败"+err.Error())
			return
		}
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
	// 获取当前用户
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}
	// 先删除历史版本记录
	for _, item := range req.Ids {
		if item == 1 && ctxUser.Username != "root" {
			response.Fail(c, nil, "删除配置失败, 通用防护策略不允许删除")
			return
		}
		raspConfig, err := r.RaspConfigRepository.GetRaspConfigById(item)
		if err != nil {
			response.Fail(c, nil, "获取配置信息失败: "+err.Error())
			return
		}
		err = r.RaspConfigHistoryRepository.DeleteRaspConfigHistory(raspConfig.RowGuid)
		if err != nil {
			response.Fail(c, nil, "删除配置历史记录信息失败: "+err.Error())
			return
		}
	}

	// 删除接口
	err = r.RaspConfigRepository.DeleteRaspConfig(req.Ids)
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
	configHistory, err := l.RaspConfigHistoryRepository.GetRaspConfigHistoryDataByGuid(config.RowGuid, config.Version)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"key": name, "value": configHistory.AgentConfigs}, "获取配置成功")
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
	hostList, _, err := hostRepository.GetRaspHostsByConfigId(req.ID)
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
	// 标记已推送
	raspConfig, err := l.RaspConfigRepository.GetRaspConfigById(req.ID)
	if err != nil {
		response.Fail(c, nil, "更新配置失败, "+err.Error())
		return
	}
	raspConfig.IsModified = false
	err = l.RaspConfigRepository.UpdateRaspConfig(raspConfig)
	if err != nil {
		response.Fail(c, nil, "更新配置失败, "+err.Error())
		return
	}
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
		Name:      srcConfig.Name + "_Copy",
		Desc:      srcConfig.Desc,
		Status:    srcConfig.Status,
		Creator:   ctxUser.Username,
		Operator:  ctxUser.Username,
		IsDefault: false,
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
		response.Fail(c, nil, "复制当前模块失败"+err.Error())
		return
	}
	var results []map[string]interface{}
	for _, module := range modules {
		var item = make(map[string]interface{})
		item["moduleName"] = module.ModuleName
		item["parameters"] = module.Parameters
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
	srcConfigHistory, err := l.RaspConfigHistoryRepository.GetRaspConfigHistoryDataByGuid(srcConfig.RowGuid, srcConfig.Version)
	if err != nil {
		response.Fail(c, nil, "获取当前配置版本失败")
		return
	}
	// 获取module信息
	var moduleConfigs []model.ModuleConfig
	var libConfigs map[string]interface{}
	var binConfigs map[string]interface{}
	var raspModules []model.RaspModule
	var raspComponents []model.RaspComponent
	err = json.Unmarshal(srcConfigHistory.ModuleConfigs, &moduleConfigs)
	if err != nil {
		response.Fail(c, nil, "导出当前配置失败")
		return
	}
	err = json.Unmarshal(srcConfigHistory.RaspBinInfo, &binConfigs)
	if err != nil {
		response.Fail(c, nil, "导出当前配置失败")
		return
	}
	err = json.Unmarshal(srcConfigHistory.RaspLibInfo, &libConfigs)
	if err != nil {
		response.Fail(c, nil, "导出当前配置失败")
		return
	}
	for _, item := range moduleConfigs {
		moduleInfo, err := l.RaspModuleRepository.GetRaspModuleByName(item.ModuleName)
		if err != nil {
			response.Fail(c, nil, "导出当前配置失败")
			return
		}
		raspModules = append(raspModules, *moduleInfo)
		componentList, _, err := l.RaspComponentRepository.GetRaspComponentsByGuid(moduleInfo.RowGuid)
		if err != nil {
			response.Fail(c, nil, "导出当前配置失败")
			return
		}
		for _, componentInfo := range componentList {
			raspComponents = append(raspComponents, *componentInfo)
		}
	}
	// 获取file信息
	var raspFiles []model.RaspFile
	for _, item := range raspComponents {
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
	if fileInfo != nil {
		raspFiles = append(raspFiles, *fileInfo)
	}
	fileInfo, err = l.RaspFileRepository.GetRaspFileByHash(binConfigs["md5"].(string))
	if err != nil {
		response.Fail(c, nil, "导出当前配置失败")
		return
	}
	if fileInfo != nil {
		raspFiles = append(raspFiles, *fileInfo)
	}
	// 构件配置对象
	var exportConfig model.RaspExportConfig
	exportConfig.RaspConfigInfo = *srcConfig
	exportConfig.RaspConfigHistoryInfo = *srcConfigHistory
	exportConfig.RaspModuleInfo = raspModules
	exportConfig.RaspComponentInfo = raspComponents
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
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fmt.Sprintf("export_%v_v%v.zip", srcConfig.Name, srcConfig.Version)))
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
		moduleInfo, err := l.RaspModuleRepository.GetRaspModuleByName(item.ModuleName)
		if err != nil {
			response.Fail(c, nil, fmt.Sprintf("查询module信息失败, %v", err))
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
				Creator:       ctxUser.Username,
				Operator:      ctxUser.Username,
			}
			err = l.RaspModuleRepository.CreateRaspModule(raspModule)
			if err != nil {
				response.Fail(c, nil, fmt.Sprintf("新建module记录失败, %v", err))
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
			moduleInfo.Creator = ctxUser.Username
			moduleInfo.Operator = ctxUser.Username
			err = l.RaspModuleRepository.UpdateRaspModule(moduleInfo)
			if err != nil {
				response.Fail(c, nil, fmt.Sprintf("更新module信息失败, %v", err))
				return
			}
		}
	}
	// 写入component信息
	for _, item := range exportConfig.RaspComponentInfo {
		componentInfo, err := l.RaspComponentRepository.GetRaspComponentsByGuidAndName(item.ParentGuid, item.ComponentName)
		if err != nil {
			response.Fail(c, nil, fmt.Sprintf("查询component信息失败, %v", err))
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
			err = l.RaspComponentRepository.CreateRaspComponent(raspComponent)
			if err != nil {
				response.Fail(c, nil, fmt.Sprintf("新建component记录失败, %v", err))
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
			err = l.RaspComponentRepository.UpdateRaspComponent(componentInfo)
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
			RowGuid:   exportConfig.RaspConfigInfo.RowGuid,
			Name:      exportConfig.RaspConfigInfo.Name,
			Version:   1,
			Desc:      exportConfig.RaspConfigInfo.Desc,
			Status:    exportConfig.RaspConfigInfo.Status,
			Creator:   ctxUser.Username,
			Operator:  ctxUser.Username,
			IsDefault: exportConfig.RaspConfigInfo.IsDefault,
		}
		err = l.RaspConfigRepository.CreateRaspConfig(raspConfig)
		if err != nil {
			response.Fail(c, nil, fmt.Sprintf("创建config信息失败, %v", err))
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
		err = l.RaspConfigHistoryRepository.CreateRaspConfigHistory(raspConfigHistory)
		if err != nil {
			response.Fail(c, nil, fmt.Sprintf("创建config版本信息失败, %v", err))
			return
		}
	} else {
		configInfo.Name = exportConfig.RaspConfigInfo.Name
		configInfo.Version += 1
		configInfo.Desc = exportConfig.RaspConfigInfo.Desc
		configInfo.Status = exportConfig.RaspConfigInfo.Status
		configInfo.Creator = ctxUser.Username
		configInfo.Operator = ctxUser.Username
		configInfo.IsDefault = exportConfig.RaspConfigInfo.IsDefault
		err = l.RaspConfigRepository.UpdateRaspConfig(configInfo)
		if err != nil {
			response.Fail(c, nil, fmt.Sprintf("更新config信息失败, %v", err))
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
		err = l.RaspConfigHistoryRepository.CreateRaspConfigHistory(raspConfigHistory)
		if err != nil {
			response.Fail(c, nil, fmt.Sprintf("创建config版本信息失败, %v", err))
			return
		}
	}
	response.Success(c, nil, "导入配置成功")
}

func (l RaspConfigController) SyncRaspConfig(c *gin.Context) {
	var req vo.SyncRaspConfigRequest
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
	srcConfig, err := l.RaspConfigRepository.GetRaspConfigById(req.SrcConfigId)
	if err != nil {
		response.Fail(c, nil, "同步策略失败，无法获取源配置信息")
		return
	}
	srcConfigHistory, err := l.RaspConfigHistoryRepository.GetRaspConfigHistoryDataByGuid(srcConfig.RowGuid, req.SrcConfigVersion)
	if err != nil {
		response.Fail(c, nil, "同步策略失败，无法获取源配置版本信息")
		return
	}
	dstConfig, err := l.RaspConfigRepository.GetRaspConfigById(req.DstConfigId)
	if err != nil {
		response.Fail(c, nil, "同步策略失败，无法获取目标配置信息")
		return
	}
	dstConfigHistory, err := l.RaspConfigHistoryRepository.GetRaspConfigHistoryDataByGuid(dstConfig.RowGuid, dstConfig.Version)
	if err != nil {
		response.Fail(c, nil, "同步策略失败，无法获取目标配置版本信息")
		return
	}
	dstConfigHistoryList, _, err := l.RaspConfigHistoryRepository.GetRaspConfigHistoryByGuid(dstConfig.RowGuid)
	if err != nil {
		response.Fail(c, nil, "同步策略失败，无法获取目标配置版本信息")
		return
	}
	var raspConfigHistory model.RaspConfigHistory
	if req.SyncOptions == 1 {
		// 仅更新高级配置信息
		newModuleConfigs, err := l.updateModuleConfig(dstConfigHistory.ModuleConfigs)
		if err != nil {
			response.Fail(c, nil, "同步策略失败，更新防护模块信息失败")
			return
		}
		raspConfigHistory = model.RaspConfigHistory{
			ParentGuid: dstConfig.RowGuid,
			Version:    dstConfigHistoryList[0].Version + 1,
			// 不改变原来的配置
			AgentMode:     dstConfigHistory.AgentMode,
			ModuleConfigs: newModuleConfigs,
			LogPath:       dstConfigHistory.LogPath,
			AgentConfigs:  dstConfigHistory.AgentConfigs,
			// 使用最新同步的数据
			RaspBinInfo: srcConfigHistory.RaspBinInfo,
			RaspLibInfo: srcConfigHistory.RaspLibInfo,
			Desc:        srcConfigHistory.Desc,
			Source:      "同步版本",
		}
	} else {
		// 覆盖更新
		raspConfigHistory = model.RaspConfigHistory{
			ParentGuid:    dstConfig.RowGuid,
			Version:       dstConfigHistoryList[0].Version + 1,
			AgentMode:     srcConfigHistory.AgentMode,
			ModuleConfigs: srcConfigHistory.ModuleConfigs,
			LogPath:       srcConfigHistory.LogPath,
			AgentConfigs:  srcConfigHistory.AgentConfigs,
			RaspBinInfo:   srcConfigHistory.RaspBinInfo,
			RaspLibInfo:   srcConfigHistory.RaspLibInfo,
			Desc:          srcConfigHistory.Desc,
			Source:        "同步版本",
		}
	}
	err = l.RaspConfigHistoryRepository.CreateRaspConfigHistory(&raspConfigHistory)
	if err != nil {
		response.Fail(c, nil, "同步配置版本失败"+err.Error())
		return
	}
	dstConfig.Version = raspConfigHistory.Version
	dstConfig.IsModified = true
	err = l.RaspConfigRepository.UpdateRaspConfig(dstConfig)
	if err != nil {
		response.Fail(c, nil, "同步配置版本失败"+err.Error())
		return
	}
	response.Success(c, nil, fmt.Sprintf("同步配置成功，新版本为V%v", raspConfigHistory.Version))
}

// 使用最新的下载地址和版本
func (l RaspConfigController) generateModuleConfig(srcModuleConfig datatypes.JSON) (datatypes.JSON, error) {
	var result []map[string]interface{}
	err := json.Unmarshal(srcModuleConfig, &result)
	if err != nil {
		return nil, err
	}
	for index, moduleConfig := range result {
		moduleInfo, err := l.RaspModuleRepository.GetRaspModuleByName(moduleConfig["moduleName"].(string))
		if err != nil {
			return nil, err
		}
		componentList, _, err := l.RaspComponentRepository.GetRaspComponentsByGuid(moduleInfo.RowGuid)
		if err != nil {
			return nil, err
		}
		var components []map[string]interface{}
		for _, item := range componentList {
			var componentInfo = make(map[string]interface{})
			componentInfo["componentName"] = item.ComponentName
			componentInfo["componentType"] = item.ComponentType
			componentInfo["componentVersion"] = item.ComponentVersion
			componentInfo["downLoadURL"] = item.DownLoadURL
			componentInfo["md5"] = item.Md5
			componentInfo["parameters"] = item.Parameters
			components = append(components, componentInfo)
		}
		result[index]["components"] = components
	}

	data, err := json.Marshal(result)
	return data, err
}

func (l RaspConfigController) updateModuleConfig(srcModuleConfig datatypes.JSON) (datatypes.JSON, error) {
	var result []model.RaspModuleConfig
	err := json.Unmarshal(srcModuleConfig, &result)
	if err != nil {
		return nil, err
	}
	// 先更新jar包的md5和下载地址
	for index, _ := range result {
		moduleConfig := &result[index]
		moduleInfo, err := l.RaspModuleRepository.GetRaspModuleByName(moduleConfig.ModuleName)
		if err != nil {
			return nil, err
		}
		if moduleInfo == nil {
			continue
		}
		componentList, _, err := l.RaspComponentRepository.GetRaspComponentsByGuid(moduleInfo.RowGuid)
		if err != nil {
			return nil, err
		}
		// 将components对象更新至新版
		moduleConfig.Components = []model.ComponentConfig{}
		for _, item := range componentList {
			moduleConfig.Components = append(moduleConfig.Components, model.ComponentConfig{
				ComponentName:    item.ComponentName,
				ComponentType:    item.ComponentType,
				ComponentVersion: item.ComponentVersion,
				DownLoadURL:      item.DownLoadURL,
				Md5:              item.Md5,
				Parameters:       item.Parameters,
			})
		}
		// 增量更新parameter对象
		var raspModuleUserConfig model.RaspModuleUserConfig
		err = json.Unmarshal(moduleInfo.Parameters, &raspModuleUserConfig)
		if err != nil {
			return nil, err
		}
		// 删除缺失过时的开关
		var delKey []string
		for k, _ := range moduleConfig.Parameters.Action {
			_, ok := raspModuleUserConfig.Action[k]
			if !ok {
				delKey = append(delKey, k)
			}
		}
		for _, k := range delKey {
			delete(moduleConfig.Parameters.Action, k)
		}
		delKey = []string{}
		for k, _ := range moduleConfig.Parameters.CnMap {
			_, ok := raspModuleUserConfig.CnMap[k]
			if !ok {
				delKey = append(delKey, k)
			}
		}
		for _, k := range delKey {
			delete(moduleConfig.Parameters.CnMap, k)
		}
		// 添加新增的action
		var addKey []string
		for k, _ := range raspModuleUserConfig.Action {
			_, ok := moduleConfig.Parameters.Action[k]
			if !ok {
				addKey = append(addKey, k)
			}
		}
		for _, k := range addKey {
			moduleConfig.Parameters.Action[k] = raspModuleUserConfig.Action[k]
		}
		addKey = []string{}
		for k, _ := range raspModuleUserConfig.CnMap {
			_, ok := moduleConfig.Parameters.CnMap[k]
			if !ok {
				addKey = append(addKey, k)
			}
		}
		for _, k := range addKey {
			moduleConfig.Parameters.CnMap[k] = raspModuleUserConfig.CnMap[k]
		}
	}
	return json.Marshal(result)
}
