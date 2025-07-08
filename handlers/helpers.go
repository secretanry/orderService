package handlers

import (
	"github.com/gin-gonic/gin"

	"wb-L0/modules/context"
)

func GetApiContext(c *gin.Context) *context.ApiContext {
	return c.MustGet("ApiContext").(*context.ApiContext)
}
