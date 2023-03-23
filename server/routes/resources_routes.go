package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"server/common"
	"server/config"
	"server/resources"
)

type HtmlHandler struct{}

func NewHtmlHandler() *HtmlHandler {
	return &HtmlHandler{}
}

// RedirectIndex 重定向
func (h *HtmlHandler) RedirectIndex(c *gin.Context) {
	c.Redirect(http.StatusFound, path.Join("/", config.Conf.System.UrlPathPrefix))
	return
}

func (h *HtmlHandler) Index(c *gin.Context) {
	c.Header("content-type", "text/html;charset=utf-8")
	c.String(200, string(resources.Html))
	return
}

func InitStaticRouter(r *gin.RouterGroup, engine *gin.Engine) gin.IRoutes {
	r.StaticFS("/static", http.FS(common.NewResource()))
	html := NewHtmlHandler()
	router := r.Group("/")
	{
		router.GET("", html.Index)
	}
	// 解决刷新404问题
	//engine.NoRoute(html.RedirectIndex)
	return r
}
