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

func (h *HtmlHandler) Logo(c *gin.Context) {
	c.Data(200, "image/svg+xml", resources.Svg)
	return
}

func (h *HtmlHandler) Favicon(c *gin.Context) {
	c.Data(200, "image/vnd.microsoft.icon", resources.Favicon)
	return
}

func InitStaticRouter(r *gin.RouterGroup, engine *gin.Engine) gin.IRoutes {
	r.StaticFS("/static", http.FS(common.NewResource()))
	r.Static("/install", "install")
	html := NewHtmlHandler()
	router := r.Group("/")
	{
		router.GET("", html.Index)
		router.GET("/webmini.svg", html.Logo)
		router.GET("/favicon.ico", html.Favicon)
	}
	// 解决刷新404问题
	//engine.NoRoute(html.RedirectIndex)
	return r
}
