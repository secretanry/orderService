package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPurchase(c *gin.Context) {
	ctx := GetApiContext(c)
	orderId, has := ctx.Params.Get("order_id")
	if !has {
		ctx.ApiError(http.StatusBadRequest, "order_id is required")
		return
	}
	ctx.JSON(200, gin.H{
		"order_id": orderId,
	})
}
