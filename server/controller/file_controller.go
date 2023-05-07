package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"server/common"
	"server/model"
	"server/repository"
	"server/response"
	"server/util"
	"server/vo"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/h2non/filetype"
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
	//pwd, _ := os.Getwd()
	form, _ := c.MultipartForm()
	files := form.File["files"]

	// 创建data目录
	dataPath := filepath.Join(DATA_DIR)
	if exist, _ := fileExist(dataPath); !exist {
		err := os.MkdirAll(dataPath, FILE_PERM)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
	}

	for _, file := range files {
		// 创建时间戳
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		// 创建目录树
		dir := filepath.Join(dataPath, file.Filename, ts)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
		filePath := filepath.Join(dir, file.Filename)
		// 保存至本地磁盘
		err = c.SaveUploadedFile(file, filePath)
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
			Timestamp:   ts,
			FileName:    file.Filename,
			FileHash:    hash,
			DiskPath:    filePath,
			DownLoadUrl: "/base/file/download?hash=" + hash,
			MimeType:    kind.MIME.Value,
			Creator:     ctxUser.Username,
		}
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
	var req vo.RaspFileDownloadRequest
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
	raspFile, err := f.RaspFileRepository.GetRaspFileByHash(req.FileHash)
	if err != nil {
		response.Fail(c, nil, "获取文件失败")
		return
	}

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", raspFile.FileName))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")

	c.File(raspFile.DiskPath)
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
	// 读取jar包manifest文件
	manifest, err := util.ReadFile(raspFile.DiskPath)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 获取默认参数配置
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
