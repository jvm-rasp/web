package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/common"
	"server/model"
	"server/repository"
	"server/response"
	"server/util"
	"server/vo"
	"time"
)

type IRaspModuleController interface {
	CreateRaspModule(c *gin.Context)
	GetRaspModules(c *gin.Context)
	UpdateRaspModules(c *gin.Context)
	BatchDeleteModuleByIds(c *gin.Context)
	DeleteModuleById(c *gin.Context)
}

type RaspModuleController struct {
	RaspConfigRepository repository.IRaspModuleRepository
}

func NewRaspModuleController() IRaspModuleController {
	raspModuleRepository := repository.NewRaspModuleRepository()
	raspModuleController := RaspModuleController{RaspConfigRepository: raspModuleRepository}
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
	raspConfigs, total, err := r.RaspConfigRepository.GetRaspModules(&req)
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

	// 获取当前用户
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}

	raspConfig := model.RaspModule{
		ModuleName:        req.ModuleName,
		ModuleVersion:     req.ModuleVersion,
		ModuleType:        req.ModuleType,
		DownLoadURL:       req.DownLoadURL,
		Md5:               req.Md5,
		MiddlewareName:    req.MiddlewareName,
		MiddlewareVersion: req.MiddlewareVersion,
		Desc:              req.Desc,
		Status:            req.Status,
		Parameters:        util.Struct2Json(req.Parameters),
		Creator:           ctxUser.Username,
		Operator:          ctxUser.Username,
		CreateTime:        time.Now().Format("2006-01-02 15:04:05"),
		UpdateTime:        time.Now().Format("2006-01-02 15:04:05"),
	}

	// 获取
	err = r.RaspConfigRepository.CreateRaspModule(&raspConfig)
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
	module, err := r.RaspConfigRepository.GetRaspModuleById(id)
	if err != nil {
		response.Fail(c, nil, "获取当前模块失败")
		return
	}

	module.ModuleName = req.ModuleName
	module.ModuleVersion = req.ModuleVersion
	module.ModuleType = req.ModuleType
	module.DownLoadURL = req.DownLoadURL
	module.Md5 = req.Md5
	module.MiddlewareName = req.MiddlewareName
	module.MiddlewareVersion = req.MiddlewareVersion
	module.Desc = req.Desc
	module.Status = req.Status
	module.Parameters = util.Struct2Json(req.Parameters)
	module.Operator = ctxUser.Username
	module.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	err = r.RaspConfigRepository.UpdateRaspModule(module)
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
	err := r.RaspConfigRepository.DeleteRaspModule(ids)
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
	err := r.RaspConfigRepository.DeleteRaspModule(req.Ids)
	if err != nil {
		response.Fail(c, nil, "删除模块失败: "+err.Error())
		return
	}
	response.Success(c, nil, "删除模块成功")
}
