package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/common"
	"server/repository"
	"server/response"
	"server/vo"
)

type IJavaProcessInfoController interface {
	GetJavaProcessInfos(c *gin.Context)
}

type JavaProcessInfoController struct {
	JavaProcessInfoRepository repository.IJavaProcessInfoRepository
}

func NewJavaProcessInfoController() IJavaProcessInfoController {
	javaProcessInfoRepository := repository.NewJavaProcessInfoRepository()
	javaProcessInfoContorller := JavaProcessInfoController{JavaProcessInfoRepository: javaProcessInfoRepository}
	return javaProcessInfoContorller
}

func (h JavaProcessInfoController) GetJavaProcessInfos(c *gin.Context) {
	var req vo.JavaProcessInfoListRequest
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
	raspHosts, total, err := h.JavaProcessInfoRepository.GetJavaProcessInfos(&req)
	if err != nil {
		response.Fail(c, nil, "获取进程列表失败")
		return
	}
	response.Success(c, gin.H{
		"data": raspHosts, "total": total,
	}, "获取进程列表失败")
}
