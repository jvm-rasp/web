package Interface

import "github.com/gin-gonic/gin"

type IRaspHostController interface {
	GetRaspHosts(c *gin.Context)
	GetRaspHost(c *gin.Context)
	BatchDeleteHostByIds(c *gin.Context)
	PushConfig(c *gin.Context)
	UpdateConfig(c *gin.Context)
	PushHostsConfig(hostList []string, content []byte) []string
	GeneratePushConfig(configId uint) ([]byte, error)
	AddHost(c *gin.Context)
	InitPushConfigService()
}
