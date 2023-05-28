package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"server/common"
	"server/model"
	"server/repository"
	"server/response"
	"server/util"
	"server/vo"
	"strings"
)

type IRaspModuleController interface {
	CreateRaspModule(c *gin.Context)
	GetRaspModules(c *gin.Context)
	UpdateRaspModules(c *gin.Context)
	BatchDeleteModuleByIds(c *gin.Context)
	DeleteModuleById(c *gin.Context)
	UpdateRaspModuleStatusById(c *gin.Context)
	UpgradeRaspModuleById(c *gin.Context)
	DeleteComponentById(c *gin.Context)
	GetRaspComponents(c *gin.Context)
}

type RaspModuleController struct {
	RaspModuleRepository    repository.IRaspModuleRepository
	RaspFileRepository      repository.IRaspFileRepository
	RaspComponentRepository repository.IRaspComponentRepository
}

func NewRaspModuleController() IRaspModuleController {
	repo1 := repository.NewRaspModuleRepository()
	repo2 := repository.NewRaspFileRepository()
	repo3 := repository.NewRaspComponentRepository()
	raspModuleController := RaspModuleController{
		RaspModuleRepository:    repo1,
		RaspFileRepository:      repo2,
		RaspComponentRepository: repo3,
	}
	return raspModuleController
}

func (r RaspModuleController) GetRaspModules(c *gin.Context) {
	var req vo.RaspModuleListRequest
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
	raspModules, total, err := r.RaspModuleRepository.GetRaspModules(&req)
	if err != nil {
		response.Fail(c, nil, "获取模块列表失败")
		return
	}
	for index, item := range raspModules {
		components, _, err := r.RaspComponentRepository.GetRaspComponentsByGuid(item.RowGuid)
		if err != nil {
			response.Fail(c, nil, "获取组件列表失败")
			return
		}
		raspModules[index].Components = components
	}

	response.Success(c, gin.H{
		"list": raspModules, "total": total,
	}, "获取模块列表成功")
}

func (r RaspModuleController) CreateRaspModule(c *gin.Context) {
	var req vo.CreateRaspModuleRequest
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
	// 判断库中是否有相同模块
	record, err := r.RaspModuleRepository.GetRaspModuleByName(req.ModuleName)
	if err != nil || record != nil {
		response.Fail(c, nil, "当前库中已存在相同模块，请重命名或者更新版本")
		return
	}
	// 获取当前用户
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}

	rowGuid := uuid.New().String()
	raspModule := model.RaspModule{
		RowGuid:       rowGuid,
		ModuleName:    req.ModuleName,
		ModuleVersion: req.ModuleVersion,
		Desc:          req.Desc,
		Parameters:    req.Parameters,
		Creator:       ctxUser.Username,
		Operator:      ctxUser.Username,
	}
	// 先处理组件列表
	for _, item := range req.Components {
		componentInfo, err := r.RaspComponentRepository.GetRaspComponentByName(item.ComponentName)
		if err != nil {
			response.Fail(c, nil, "查询防护组件失败"+err.Error())
			return
		}
		if componentInfo != nil {
			response.Fail(c, nil, fmt.Sprintf("创建防护模块失败: 组件%v已经被引用无法重复引用", componentInfo.ComponentName))
			return
		}
		raspComponent := model.RaspComponent{
			ParentGuid:       rowGuid,
			ComponentName:    item.ComponentName,
			ComponentType:    item.ComponentType,
			ComponentVersion: item.ComponentVersion,
			DownLoadURL:      item.DownLoadURL,
			Md5:              item.Md5,
			Parameters:       item.Parameters,
		}
		err = r.RaspComponentRepository.CreateRaspComponent(&raspComponent)
		if err != nil {
			response.Fail(c, nil, "创建防护组件失败"+err.Error())
			return
		}
	}
	// 创建防护模块
	err = r.RaspModuleRepository.CreateRaspModule(&raspModule)
	if err != nil {
		response.Fail(c, nil, "创建防护模块失败"+err.Error())
		return
	}
	response.Success(c, nil, "创建防护模块成功")
	return
}

func (r RaspModuleController) UpdateRaspModules(c *gin.Context) {
	var req vo.UpdateRaspModuleRequest
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
	module, err := r.RaspModuleRepository.GetRaspModuleById(id)
	if err != nil {
		response.Fail(c, nil, "获取当前模块失败")
		return
	}
	module.ModuleName = req.ModuleName
	module.ModuleVersion = req.ModuleVersion
	module.Parameters = req.Parameters
	module.Desc = req.Desc
	module.Creator = ctxUser.Username
	module.Operator = ctxUser.Username
	err = r.RaspModuleRepository.UpdateRaspModule(module)
	if err != nil {
		response.Fail(c, nil, "更新当前模块失败")
		return
	}
	err = r.RaspComponentRepository.DeleteRaspComponentByGuid(module.RowGuid)
	if err != nil {
		response.Fail(c, nil, "更新当前模块失败")
		return
	}
	for _, item := range req.Components {
		raspComponent := model.RaspComponent{
			ParentGuid:       module.RowGuid,
			ComponentName:    item.ComponentName,
			ComponentType:    item.ComponentType,
			ComponentVersion: item.ComponentVersion,
			DownLoadURL:      item.DownLoadURL,
			Md5:              item.Md5,
			Parameters:       item.Parameters,
		}
		err = r.RaspComponentRepository.CreateRaspComponent(&raspComponent)
		if err != nil {
			response.Fail(c, nil, "创建防护组件失败"+err.Error())
			return
		}
	}
	response.Success(c, nil, "更新模块成功")
}

func (r RaspModuleController) DeleteModuleById(c *gin.Context) {
	var req vo.DeleteRaspModuleRequest
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
	ids := []uint{req.Id}
	var rowGuids []string
	for _, id := range ids {
		module, err := r.RaspModuleRepository.GetRaspModuleById(id)
		if err != nil {
			continue
		}
		rowGuids = append(rowGuids, module.RowGuid)
	}
	// 删除接口
	err := r.RaspModuleRepository.DeleteRaspModule(ids)
	if err != nil {
		response.Fail(c, nil, "删除模块失败: "+err.Error())
		return
	}
	err = r.RaspComponentRepository.DeleteRaspComponentByGuids(rowGuids)
	if err != nil {
		response.Fail(c, nil, "删除组件失败: "+err.Error())
		return
	}
	response.Success(c, nil, "删除模块成功")
}

// 批量删除接口
func (r RaspModuleController) BatchDeleteModuleByIds(c *gin.Context) {
	var req vo.DeleteBatchRaspModuleRequest
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
	var rowGuids []string
	for _, id := range req.Ids {
		module, err := r.RaspModuleRepository.GetRaspModuleById(id)
		if err != nil {
			continue
		}
		rowGuids = append(rowGuids, module.RowGuid)
	}
	// 删除接口
	err := r.RaspModuleRepository.DeleteRaspModule(req.Ids)
	if err != nil {
		response.Fail(c, nil, "删除模块失败: "+err.Error())
		return
	}
	err = r.RaspComponentRepository.DeleteRaspComponentByGuids(rowGuids)
	if err != nil {
		response.Fail(c, nil, "删除组件失败: "+err.Error())
		return
	}
	response.Success(c, nil, "删除模块成功")
}

func (r RaspModuleController) UpdateRaspModuleStatusById(c *gin.Context) {
	var req vo.UpdateRaspModuleStatusRequest
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
	module, err := r.RaspModuleRepository.GetRaspModuleById(id)
	if err != nil {
		response.Fail(c, nil, "获取当前模块失败")
		return
	}
	module.Operator = ctxUser.Username

	err = r.RaspModuleRepository.UpdateRaspModule(module)
	if err != nil {
		response.Fail(c, nil, "更新当前模块失败")
		return
	}
	response.Success(c, nil, "更新模块成功")
}

func (r RaspModuleController) UpgradeRaspModuleById(c *gin.Context) {
	var req vo.UpgradeRaspModuleRequest
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
	module, err := r.RaspModuleRepository.GetRaspModuleById(req.ID)
	if err != nil {
		response.Fail(c, nil, "获取当前模块信息失败")
		return
	}
	if module.Upgradable {
		upgradeComponents, _, err := r.RaspComponentRepository.GetRaspUpgradeComponentsByGuid(module.RowGuid)
		if err != nil {
			response.Fail(c, nil, "获取最新模块信息失败")
			return
		}
		for _, item := range upgradeComponents {
			raspFile, err := r.RaspFileRepository.GetRaspFileByHash(item.Md5)
			if err != nil {
				response.Fail(c, nil, "获取最新模块信息失败")
				return
			}
			// 获取jar包的信息
			manifest, err := util.ReadFile(raspFile.DiskPath)
			if err != nil {
				response.Fail(c, nil, "读取jar包信息失败: "+err.Error())
				return
			}
			// 获取默认参数配置
			parametersStr, err := util.GetDefaultParameters(raspFile.DiskPath)
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}
			var defaultParameter map[string]interface{}
			err = json.Unmarshal([]byte(parametersStr), &defaultParameter)
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}
			var parameters datatypes.JSON
			parameters, err = json.Marshal(defaultParameter["parameters"])
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}
			// 更新component数据
			componentName := manifest["ModuleName"]
			componentVersion := manifest["ModuleVersion"]
			componentInfo, err := r.RaspComponentRepository.GetRaspComponentsByGuidAndName(module.RowGuid, componentName)
			if err != nil {
				response.Fail(c, nil, "获取组件信息失败, err: "+err.Error())
				return
			}
			componentInfo.ComponentName = componentName
			componentInfo.ComponentVersion = componentVersion
			componentInfo.DownLoadURL = raspFile.DownLoadUrl
			componentInfo.Md5 = raspFile.FileHash
			componentInfo.Parameters = parameters
			err = r.RaspComponentRepository.UpdateRaspComponent(componentInfo)
			if err != nil {
				response.Fail(c, nil, "更新组件信息失败, err: "+err.Error())
				return
			}
			// 如果是算法逻辑类jar包则替换参数
			if strings.HasSuffix(componentName, "-algorithm") {
				var moduleParameters = make(map[string]interface{})
				var action = make(map[string]interface{})
				moduleParameters["cn_map"] = defaultParameter["cn_map"]
				for k, v := range defaultParameter["parameters"].(map[string]interface{}) {
					if strings.HasSuffix(k, "_action") {
						action[k] = v
					}
				}
				moduleParameters["action"] = action
				module.Parameters, _ = json.Marshal(moduleParameters)
			}
			err = r.RaspComponentRepository.DeleteRaspUpgradeComponentById(item.ID)
			if err != nil {
				response.Fail(c, nil, "删除更新临时组件信息失败, err: "+err.Error())
				return
			}
		}
		// 取消标记
		module.Upgradable = false
		module.Operator = ctxUser.Username
		err = r.RaspModuleRepository.UpdateRaspModule(module)
		if err != nil {
			response.Fail(c, nil, "标记模块已更新失败")
			return
		}
	} else {
		response.Fail(c, nil, "当前模块已是最新")
		return
	}
	response.Success(c, nil, "更新模块成功, 更新配置版本后生效")
}

func (r RaspModuleController) DeleteComponentById(c *gin.Context) {
	var req vo.RaspComponentDeleteRequest
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
	ids := []uint{req.Id}
	if err := r.RaspComponentRepository.DeleteRaspComponentByIds(ids); err != nil {
		response.Fail(c, nil, "删除组件失败: "+err.Error())
		return
	}
	response.Success(c, nil, "删除模块成功")
}

func (r RaspModuleController) GetRaspComponents(c *gin.Context) {
	var req vo.RaspComponentListRequest
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
	raspComponents, total, err := r.RaspComponentRepository.GetRaspComponentsByGuid(req.ParentGuid)
	if err != nil {
		response.Fail(c, nil, "获取组件列表失败")
		return
	}
	response.Success(c, gin.H{
		"list": raspComponents, "total": total,
	}, "获取组件列表成功")
}
