package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wb-L0/services/composer/orders"
	"wb-L0/services/database"
)

// GetPurchase
// @Tags purchases
// @Summary Get order by uid
// @ID get-order-by-id
// @Param uid path string true "Order uid"
// @Produce json
// @Success 200 {object} structs.Order "Order obtained"
// @Success 404 {object} structs.ApiError "Order not found"
// @Failure 500 {object} structs.ApiError "Internal server error"
// @Router /api/orders/{uid} [get]
func GetPurchase(c *gin.Context) {
	ctx := GetApiContext(c)
	orderId, has := ctx.Params.Get("order_id")
	if !has {
		ctx.ApiError(http.StatusBadRequest, "order_id is required")
		return
	}
	order, err := orders.GetOrderById(ctx, orderId)
	if err != nil {
		if database.IsErrOrderNotFound(err) {
			ctx.ApiError(http.StatusNotFound, err.Error())
			return
		}
		ctx.ApiError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, order)
}
