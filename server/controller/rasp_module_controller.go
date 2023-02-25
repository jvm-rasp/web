package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/common"
	"server/model"
	"server/repository"
	"server/response"
	"server/vo"
)

type IRaspModuleController interface {
	CreateRaspModule(c *gin.Context)
	GetRaspModules(c *gin.Context)
	BatchDeleteModuleByIds(c *gin.Context)
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
		"data": raspConfigs, "total": total,
	}, "获取模块列表成功")
}

func (r RaspModuleController) CreateRaspModule(c *gin.Context) {
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

	raspConfig := model.RaspModule{
		Name:    req.Name,
		Tag:     req.Tag,
		Desc:    req.Desc,
		Creator: ctxUser.Username,
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

// 批量删除接口
func (r RaspModuleController) BatchDeleteModuleByIds(c *gin.Context) {
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

	// 删除接口
	err := r.RaspConfigRepository.DeleteRaspModule(req.Ids)
	if err != nil {
		response.Fail(c, nil, "删除模块失败: "+err.Error())
		return
	}
	response.Success(c, nil, "删除模块成功")
}



