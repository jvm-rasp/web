package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/h2non/filetype"
	"os"
	"path/filepath"
	"server/common"
	"server/config"
	"server/model"
	"server/repository"
	"server/response"
	"server/util"
	"server/vo"
)

type IFileController interface {
	Upload(c *gin.Context)
	Download(c *gin.Context)
	GetRaspFiles(c *gin.Context)
	Delete(c *gin.Context)
	GetModuleInfo(c *gin.Context)
}

// DATA_DIR jar包存方路径
const DATA_DIR = "upload"

const FILE_PERM = 0755

type FileController struct {
	RaspFileRepository repository.IRaspFileRepository
}

func NewFileController() IFileController {
	repository := repository.NewRaspFileRepository()
	fileController := FileController{RaspFileRepository: repository}
	return fileController
}

func (f FileController) Upload(c *gin.Context) {
	pwd, _ := os.Getwd()
	form, _ := c.MultipartForm()
	files := form.File["files"]

	// 创建data目录
	dataPath := filepath.Join(pwd, DATA_DIR)
	if exist, _ := fileExist(dataPath); !exist {
		err := os.MkdirAll(dataPath, FILE_PERM)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
	}

	for _, file := range files {
		filePath := filepath.Join(dataPath, file.Filename)
		err := c.SaveUploadedFile(file, filePath)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
		// 判断文件类型
		buf, _ := os.ReadFile(filePath)
		kind, _ := filetype.Match(buf)
		hash, err := util.GetFileMd5(filePath)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
		// 获取当前用户
		ur := repository.NewUserRepository()
		ctxUser, err := ur.GetCurrentUser(c)
		if err != nil {
			response.Fail(c, nil, "获取当前用户信息失败")
			return
		}

		fileInfo := &model.RaspFile{
			FileName:    file.Filename,
			FileHash:    hash,
			DiskPath:    filePath,
			DownLoadUrl: "/" + config.Conf.System.UrlPathPrefix + "/file/download?file=" + file.Filename,
			MimeType:    kind.MIME.Value,
			Creator:     ctxUser.Username,
		}

		// TODO 如果存在则更新
		err = f.RaspFileRepository.CreateRaspFile(fileInfo)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
	}
	response.Success(c, nil, "uploaded file success")
}

func (r FileController) GetRaspFiles(c *gin.Context) {
	var req vo.RaspFileListRequest
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
	raspFiles, total, err := r.RaspFileRepository.GetRaspFiles(&req)
	if err != nil {
		response.Fail(c, nil, "获取jar包列表失败")
		return
	}
	response.Success(c, gin.H{
		"list": raspFiles, "total": total,
	}, "获取jar包列表成功")
}

func (f FileController) Download(c *gin.Context) {
	// TODO 接口不鉴权
}

func (f FileController) Delete(c *gin.Context) {
	var req vo.RaspFileDeleteRequest
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
	// 删除附件
	var list []*model.RaspFile
	err := common.DB.Where("id IN (?)", req.Ids).Unscoped().Find(&list).Error
	if err != nil {
		response.Fail(c, nil, "删除文件失败")
		return
	}
	for _, item := range list {
		if err := os.Remove(item.DiskPath); err != nil {
			response.Fail(c, nil, err.Error())
		}
	}
	// 删除数据库
	err = f.RaspFileRepository.DeleteRaspFile(req.Ids)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "删除文件成功")
}

func (f FileController) GetModuleInfo(c *gin.Context) {
	var req vo.RaspFileInfoRequest
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
	// 获取文件路径
	raspFile, err := f.RaspFileRepository.GetRaspFileById(req.Id)
	if err != nil {
		response.Fail(c, nil, "获取jar包信息失败")
		return
	}
	// 读取jar包mainfest文件
	manifest, err := util.ReadFile(raspFile.DiskPath)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	defaultParameters, err := util.GetDefaultParameters(raspFile.DiskPath)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"manifest": manifest, "parameters": defaultParameters}, "读取文件信息成功")
}

func fileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, nil
}
