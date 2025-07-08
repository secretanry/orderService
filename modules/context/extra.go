package context

import (
	"github.com/gin-gonic/gin"

	"wb-L0/structs"
)

type ApiContext struct {
	*gin.Context
}

func (ctx *ApiContext) ApiError(code int, message string) {
	ctx.JSON(code, structs.ApiError{
		Message: message,
	})
}
