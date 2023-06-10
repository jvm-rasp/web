package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gookit/goutil/fsutil"
	"github.com/imroc/req"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	url2 "net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"server/common"
	"server/config"
	"server/model"
	"server/report"
	"server/repository"
	"server/response"
	"server/vo"
	"strconv"
)

type ISystemSettingController interface {
	Update(c *gin.Context)
	List(c *gin.Context)
	GetProjectInfo(c *gin.Context)
}

type SystemSettingController struct {
	SystemSettingRepository repository.ISystemSettingRepository
}

func NewSystemSettingController() ISystemSettingController {
	systemSettingController := SystemSettingController{
		SystemSettingRepository: repository.NewSystemSettingRepository(),
	}
	return systemSettingController
}

func (this SystemSettingController) Update(c *gin.Context) {
	var req vo.UpdateSettingRequest
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
	setting, err := this.SystemSettingRepository.GetSettingByName(req.Name)
	if err != nil {
		response.Fail(c, nil, fmt.Sprintf("更新配置失败, error: %v", err))
		return
	}
	if setting != nil {
		setting.Value = AnyToStr(req.Value)
		if err = this.SystemSettingRepository.SaveSetting(setting); err != nil {
			response.Fail(c, nil, fmt.Sprintf("更新配置失败, error: %v", err))
		}
		// 判断如果是连接更新服务器
		switch setting.Name {
		case "autoUpdate":
			err = HandleAutoUpdate(setting.Value)
		case "reportUrl":
			err = HandleReportUrl(setting.Value)
		}
		if err != nil {
			response.Fail(c, nil, fmt.Sprintf("更新配置失败: %v", err.Error()))
			return
		}
		response.Success(c, nil, "更新配置成功")
	} else {
		response.Fail(c, nil, fmt.Sprintf("未找到配置: %v", req.Name))
	}
}

func (this SystemSettingController) List(c *gin.Context) {
	settings, err := this.SystemSettingRepository.GetSettings()
	if err != nil {
		response.Fail(c, nil, "获取系统设置列表失败")
		return
	}
	var result = make(map[string]interface{})
	for _, item := range settings {
		switch item.Type {
		case reflect.String.String():
			result[item.Name] = item.Value
		case reflect.Bool.String():
			result[item.Name], _ = strconv.ParseBool(item.Value)
		case reflect.Int.String():
			result[item.Name], _ = strconv.Atoi(item.Value)
		}
	}
	response.Success(c, gin.H{
		"list": result,
	}, "获取配置列表成功")
}

func (this SystemSettingController) GetProjectInfo(c *gin.Context) {
	var request vo.GetProjectInfoRequest
	// 参数绑定
	if err := c.ShouldBind(&request); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&request); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	url, err := url2.Parse(request.ReportUrl)
	if err != nil {
		response.Fail(c, nil, "获取项目guid失败, 请检查reportUrl是否正确: "+err.Error())
		return
	}
	urlPath := path.Join(url.Path, "rest", "getProjectInfo")
	resp, err := req.Get(fmt.Sprintf("%v://%v%v", url.Scheme, url.Host, urlPath))
	if err != nil {
		response.Fail(c, nil, "获取项目guid失败, 请检查reportUrl是否正确: "+err.Error())
		return
	}
	var projectInfo vo.ProjectInfo
	err = resp.ToJSON(&projectInfo)
	if err != nil {
		response.Fail(c, nil, "获取项目guid失败, 请检查reportUrl是否正确: "+err.Error())
		return
	}
	setting, err := this.SystemSettingRepository.GetSettingByName("projectGuid")
	if err != nil {
		response.Fail(c, nil, "更新projectGuid配置项失败: "+err.Error())
		return
	}
	if setting != nil {
		setting.Value = projectInfo.ProjectGuid
		if err = this.SystemSettingRepository.SaveSetting(setting); err != nil {
			response.Fail(c, nil, "更新projectGuid配置项失败: "+err.Error())
			return
		}
	}
	response.Success(c, nil, "获取项目guid成功")
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func AnyToStr(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func HandleAutoUpdate(value string) error {
	if value == "false" {
		report.UpdateManager.AutoUpdate = false
		report.UpdateManager.DisConnect()
	} else {
		report.UpdateManager.AutoUpdate = true
	}
	configPath := filepath.Join(config.Conf.Env.WorkDir, "filebeat", "filebeat.yml")
	viperConfig, err := readFilebeatConfig(configPath)
	if err != nil {
		return err
	}
	var filebeatConfig model.FilebeatConfig
	err = viperConfig.Unmarshal(&filebeatConfig)
	for index, _ := range filebeatConfig.Filebeat.Inputs {
		item := filebeatConfig.Filebeat.Inputs[index]
		item.Enabled = report.UpdateManager.AutoUpdate
	}
	err = writeFilebeatConfig(configPath, structs.Map(filebeatConfig))
	return err
}

func HandleReportUrl(remoteHost string) error {
	url, _ := url2.Parse(remoteHost)
	urlPath := path.Join(url.Path, "rest", "exchangedata")
	if url.Scheme == "https" {
		remoteHost = fmt.Sprintf("%v://%v%v", "https", url.Host, urlPath)
	}
	if url.Scheme == "http" {
		remoteHost = fmt.Sprintf("%v://%v%v", "http", url.Host, urlPath)
	}
	err := writeFilebeatConfigByKey("output.http.hosts", []string{remoteHost})
	return err
}

func readFilebeatConfig(fileName string) (*viper.Viper, error) {
	// 修改filebeat配置
	if !fsutil.PathExist(fileName) {
		return nil, errors.New(fmt.Sprintf("%v does not exists", fileName))
	}
	var filebeatConfig = viper.New()
	filebeatConfig.SetConfigFile(fileName)
	filebeatConfig.SetConfigType("yaml")
	if err := filebeatConfig.ReadInConfig(); err != nil {
		return nil, errors.New(fmt.Sprintf("Read Filebeat Config File Failed, error: %v", err))
	}
	return filebeatConfig, nil
}

func writeFilebeatConfig(fileName string, setting map[string]interface{}) error {
	configData, err := yaml.Marshal(setting)
	if err != nil {
		return errors.New(fmt.Sprintf("Marshal Filebeat Config File Failed, error: %v", err))
	}
	err = os.WriteFile(fileName, configData, 0755)
	if err != nil {
		return errors.New(fmt.Sprintf("Write Filebeat Config File Failed error: %v", err))
	}
	return nil
}

func writeFilebeatConfigByKey(configName string, configValue interface{}) error {
	// 修改filebeat配置
	fileName := filepath.Join(config.Conf.Env.WorkDir, "filebeat", "filebeat.yml")
	if !fsutil.PathExist(fileName) {
		return errors.New(fmt.Sprintf("%v does not exists", fileName))
	}
	var filebeatConfig = viper.New()
	filebeatConfig.SetConfigFile(fileName)
	filebeatConfig.SetConfigType("yaml")
	if err := filebeatConfig.ReadInConfig(); err != nil {
		return errors.New(fmt.Sprintf("Read Filebeat Config File Failed, error: %v", err))
	}
	filebeatConfig.Set(configName, configValue)
	setting := filebeatConfig.AllSettings()
	configData, err := yaml.Marshal(setting)
	if err != nil {
		return errors.New(fmt.Sprintf("Marshal Filebeat Config File Failed, error: %v", err))
	}
	err = os.WriteFile(fileName, configData, 0755)
	if err != nil {
		return errors.New(fmt.Sprintf("Write Filebeat Config File Failed error: %v", err))
	}
	return nil
}
