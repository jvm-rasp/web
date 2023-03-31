package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"server/response"
)

type IFileController interface {
	Upload(c *gin.Context)
	Download(c *gin.Context)
}

// DATA_DIR jar包存方路径
const DATA_DIR = "data"

// FileMap 存储jar包信息
var FileMap = make(map[string]FileInfo)

type FileInfo struct {
	FileName      string
	FileHash      string
	DiskPath      string
	DownLoadUrl   string
	ModuleName    string
	ModuleVersion string
	UpdateTime    string
}

type FileController struct {
}

func NewFileController() IFileController {
	fileController := FileController{}
	return fileController
}

func (f FileController) Upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	pwd, _ := os.Getwd()
	err := c.SaveUploadedFile(file, filepath.Join(pwd, DATA_DIR, file.Filename))
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func (f FileController) Download(c *gin.Context) {
	//
}
