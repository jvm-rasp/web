package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	UpGradeRaspModuleById(c *gin.Context)
}

type RaspModuleController struct {
	RaspModuleRepository repository.IRaspModuleRepository
	RaspFileRepository   repository.IRaspFileRepository
}

func NewRaspModuleController() IRaspModuleController {
	repo1 := repository.NewRaspModuleRepository()
	repo2 := repository.NewRaspFileRepository()
	raspModuleController := RaspModuleController{
		RaspModuleRepository: repo1,
		RaspFileRepository:   repo2,
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
	raspConfigs, total, err := r.RaspModuleRepository.GetRaspModules(&req)
	if err != nil {
		response.Fail(c, nil, "获取模块列表失败")
		return
	}
	response.Success(c, gin.H{
		"list": raspConfigs, "total": total,
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
	record, err := r.RaspModuleRepository.GetRaspModuleByName(req.ModuleName, req.ModuleVersion)
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

	raspConfig := model.RaspModule{
		ModuleName:    req.ModuleName,
		ModuleVersion: req.ModuleVersion,
		ModuleType:    req.ModuleType,
		DownLoadURL:   req.DownLoadURL,
		Md5:           req.Md5,
		Desc:          req.Desc,
		Status:        req.Status,
		Parameters:    req.Parameters,
		Creator:       ctxUser.Username,
		Operator:      ctxUser.Username,
	}

	// 获取
	err = r.RaspModuleRepository.CreateRaspModule(&raspConfig)
	if err != nil {
		response.Fail(c, nil, "创建模块列表失败"+err.Error())
		return
	}
	response.Success(c, nil, "创建模块成功")
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
	module.ModuleType = req.ModuleType
	module.DownLoadURL = req.DownLoadURL
	module.Md5 = req.Md5
	module.Desc = req.Desc
	module.Status = req.Status
	module.Parameters = req.Parameters
	module.Operator = ctxUser.Username

	err = r.RaspModuleRepository.UpdateRaspModule(module)
	if err != nil {
		response.Fail(c, nil, "更新当前模块失败")
		return
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
	// 删除接口
	err := r.RaspModuleRepository.DeleteRaspModule(ids)
	if err != nil {
		response.Fail(c, nil, "删除模块失败: "+err.Error())
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

	// 删除接口
	err := r.RaspModuleRepository.DeleteRaspModule(req.Ids)
	if err != nil {
		response.Fail(c, nil, "删除模块失败: "+err.Error())
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
	module.Status = !module.Status
	module.Operator = ctxUser.Username

	err = r.RaspModuleRepository.UpdateRaspModule(module)
	if err != nil {
		response.Fail(c, nil, "更新当前模块失败")
		return
	}
	response.Success(c, nil, "更新模块成功")
}

func (r RaspModuleController) UpGradeRaspModuleById(c *gin.Context) {
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
	if module.Upgradable && module.NewMd5 != "" {
		// TODO 更新为最新模块
		raspFile, err := r.RaspFileRepository.GetRaspFileByHash(module.NewMd5)
		if err != nil {
			response.Fail(c, nil, "获取最新模块信息失败")
			return
		}
		// 获取默认参数配置
		parametersStr, err := util.GetDefaultParameters(raspFile.DiskPath)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
		parametersStr = strings.ReplaceAll(parametersStr, "\n", "")
		var parameters datatypes.JSON
		err = json.Unmarshal([]byte(parametersStr), &parameters)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
		module.DownLoadURL = raspFile.DownLoadUrl
		module.Md5 = raspFile.FileHash
		module.Parameters = parameters
		// 取消标记
		module.Upgradable = false
		module.NewMd5 = ""
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
	response.Success(c, nil, "更新模块成功")
}
