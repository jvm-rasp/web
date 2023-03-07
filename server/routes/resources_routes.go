package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/common"
	"server/resources"
)

type HtmlHandler struct{}

func NewHtmlHandler() *HtmlHandler {
	return &HtmlHandler{}
}

// RedirectIndex 重定向
func (h *HtmlHandler) RedirectIndex(c *gin.Context) {
	c.Redirect(http.StatusFound, "/ui")
	return
}

func (h *HtmlHandler) Index(c *gin.Context) {
	c.Header("content-type", "text/html;charset=utf-8")
	c.String(200, string(resources.Html))
	return
}

// InitStaticRouter 静态资源配置
func InitStaticRouter(engine *gin.Engine) {
	engine.StaticFS("/static", http.FS(common.NewResource()))

	html := NewHtmlHandler()
	group := engine.Group("/ui")
	{
		group.GET("", html.Index)
	}
	// 解决刷新404问题
	engine.NoRoute(html.RedirectIndex)
}
