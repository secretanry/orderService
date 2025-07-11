package routing

import (
	"github.com/gin-gonic/gin"

	"wb-L0/handlers"
	"wb-L0/modules/context"
)

func ApiContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ApiContext", &context.ApiContext{Context: c})
		c.Next()
	}
}

func MountPurchasesRoutes(r *gin.Engine) {
	api := r.Group("/api")
	order := api.Group("/order")
	order.GET("/:order_id", handlers.GetPurchase)
}
