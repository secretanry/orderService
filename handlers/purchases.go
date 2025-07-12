package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"wb-L0/modules/monitoring"
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
	logger := monitoring.LogWithContext(c)

	orderId, has := ctx.Params.Get("order_id")
	if !has {
		logger.Warn("Missing order_id parameter")
		ctx.ApiError(http.StatusBadRequest, "order_id is required")
		return
	}

	logger.Info("Processing order retrieval request",
		zap.String("order_id", orderId))

	order, err := orders.GetOrderById(ctx, orderId)
	if err != nil {
		if database.IsErrOrderNotFound(err) {
			logger.Info("Order not found",
				zap.String("order_id", orderId))
			ctx.ApiError(http.StatusNotFound, err.Error())
			return
		}
		logger.Error("Failed to retrieve order",
			zap.String("order_id", orderId), zap.Error(err))
		ctx.ApiError(http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("Order retrieved successfully",
		zap.String("order_id", orderId))
	ctx.JSON(http.StatusOK, order)
}
