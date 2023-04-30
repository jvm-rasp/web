package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/goutil/arrutil"
	"path"
)

var staticExts = []string{".js", ".css", ".svg", ".jpg"}
var cacheAge = 1 * 24 * 3600 // 单位秒

func CacheControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if arrutil.Contains(staticExts, path.Ext(c.Request.URL.Path)) {
			c.Header("Cache-Control", fmt.Sprintf("public, max-age=%v", cacheAge))
		}
	}
}
