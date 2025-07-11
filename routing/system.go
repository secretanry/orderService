package routing

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"

	"wb-L0/handlers"
)

func MountSystemRoutes(r *gin.Engine) {
	r.Use(ApiContextMiddleware())
	r.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
		httpSwagger.DocExpansion("none"),
	)))
	r.GET("/docs/swagger.json", func(c *gin.Context) {
		ctx := handlers.GetApiContext(c)
		data, err := os.ReadFile("docs/swagger.json")
		if err != nil {
			ctx.ApiError(http.StatusInternalServerError, "Failed to read Swagger JSON")
			return
		}
		ctx.String(http.StatusOK, string(data))
	})
}
