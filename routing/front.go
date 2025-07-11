package routing

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MountFrontRoutes(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}
