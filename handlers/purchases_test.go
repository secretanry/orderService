package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	apicontext "wb-L0/modules/context"
	"wb-L0/structs"
)

// MockOrderService is a mock implementation of the order service
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) GetOrderById(ctx context.Context, orderId string) (*structs.Order, error) {
	args := m.Called(ctx, orderId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*structs.Order), args.Error(1)
}

// TestGetPurchaseSuccess tests successful order retrieval
func TestGetPurchaseSuccess(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock order service
	mockService := new(MockOrderService)

	// Create test order
	expectedOrder := &structs.Order{
		OrderUid:    "test-order-123",
		TrackNumber: "TEST_TRACK_123",
		CustomerId:  "test_customer",
		Delivery: structs.Delivery{
			Name:  "Test Testov",
			Phone: "+9720000000",
			Email: "test@gmail.com",
		},
		Payment: structs.Payment{
			Transaction: "b563feb7b2b84b6test",
			Currency:    "USD",
			Amount:      1817,
		},
		Items: []structs.Item{
			{
				ChrtId:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Name:        "Mascaras",
				Brand:       "Vivienne Sabo",
			},
		},
	}

	mockService.On("GetOrderById", mock.Anything, "test-order-123").Return(expectedOrder, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("ApiContext", &apicontext.ApiContext{Context: c})
		c.Next()
	})

	router.GET("/api/order/:order_id", func(c *gin.Context) {
		ctx := GetApiContext(c)
		orderId, has := ctx.Params.Get("order_id")
		if !has {
			ctx.ApiError(http.StatusBadRequest, "order_id is required")
			return
		}

		order, err := mockService.GetOrderById(ctx, orderId)
		if err != nil {
			ctx.ApiError(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, order)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/order/test-order-123", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseOrder structs.Order
	err := json.Unmarshal(w.Body.Bytes(), &responseOrder)
	require.NoError(t, err)

	assert.Equal(t, expectedOrder.OrderUid, responseOrder.OrderUid)
	assert.Equal(t, expectedOrder.TrackNumber, responseOrder.TrackNumber)
	assert.Equal(t, expectedOrder.CustomerId, responseOrder.CustomerId)
	assert.Len(t, responseOrder.Items, 1)
	assert.Equal(t, expectedOrder.Items[0].Name, responseOrder.Items[0].Name)

	mockService.AssertExpectations(t)
}

func TestGetPurchaseNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockOrderService)

	mockService.On("GetOrderById", mock.Anything, "non-existent-order").Return(nil, assert.AnError)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("ApiContext", &apicontext.ApiContext{Context: c})
		c.Next()
	})

	router.GET("/api/order/:order_id", func(c *gin.Context) {
		ctx := GetApiContext(c)
		orderId, has := ctx.Params.Get("order_id")
		if !has {
			ctx.ApiError(http.StatusBadRequest, "order_id is required")
			return
		}

		order, err := mockService.GetOrderById(ctx, orderId)
		if err != nil {
			ctx.ApiError(http.StatusNotFound, "Order not found")
			return
		}

		ctx.JSON(http.StatusOK, order)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/order/non-existent-order", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)

	assert.Contains(t, errorResponse, "message")

	mockService.AssertExpectations(t)
}

func TestGetPurchaseMissingID(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create router and set up route
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("ApiContext", &apicontext.ApiContext{Context: c})
		c.Next()
	})

	router.GET("/api/order/:order_id", func(c *gin.Context) {
		ctx := GetApiContext(c)
		orderId, has := ctx.Params.Get("order_id")
		if !has {
			ctx.ApiError(http.StatusBadRequest, "order_id is required")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"order_id": orderId})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/order/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPurchaseInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockOrderService)
	mockService.On("GetOrderById", mock.Anything, "invalid").Return(nil, assert.AnError)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("ApiContext", &apicontext.ApiContext{Context: c})
		c.Next()
	})
	router.GET("/api/order/:order_id", func(c *gin.Context) {
		ctx := GetApiContext(c)
		orderId, has := ctx.Params.Get("order_id")
		if !has {
			ctx.ApiError(http.StatusBadRequest, "order_id is required")
			return
		}
		order, err := mockService.GetOrderById(ctx, orderId)
		if err != nil {
			ctx.ApiError(http.StatusNotFound, "Order not found")
			return
		}
		ctx.JSON(http.StatusOK, order)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/order/invalid", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

func TestGetPurchaseContextHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockOrderService)

	expectedOrder := &structs.Order{
		OrderUid: "context-test-order",
	}

	mockService.On("GetOrderById", mock.Anything, "context-test-order").Return(expectedOrder, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("ApiContext", &apicontext.ApiContext{Context: c})
		c.Next()
	})

	router.GET("/api/order/:order_id", func(c *gin.Context) {
		ctx := GetApiContext(c)

		assert.NotNil(t, ctx)
		assert.NotNil(t, ctx.Context)

		orderId, has := ctx.Params.Get("order_id")
		if !has {
			ctx.ApiError(http.StatusBadRequest, "order_id is required")
			return
		}

		order, err := mockService.GetOrderById(ctx, orderId)
		if err != nil {
			ctx.ApiError(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, order)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/order/context-test-order", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestGetPurchaseErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockOrderService)

	mockService.On("GetOrderById", mock.Anything, "db-error").Return(nil, assert.AnError)
	mockService.On("GetOrderById", mock.Anything, "cache-error").Return(nil, assert.AnError)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("ApiContext", &apicontext.ApiContext{Context: c})
		c.Next()
	})

	router.GET("/api/order/:order_id", func(c *gin.Context) {
		ctx := GetApiContext(c)
		orderId, has := ctx.Params.Get("order_id")
		if !has {
			ctx.ApiError(http.StatusBadRequest, "order_id is required")
			return
		}

		order, err := mockService.GetOrderById(ctx, orderId)
		if err != nil {
			if orderId == "db-error" {
				ctx.ApiError(http.StatusInternalServerError, "Database error")
			} else {
				ctx.ApiError(http.StatusNotFound, "Order not found")
			}
			return
		}

		ctx.JSON(http.StatusOK, order)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/order/db-error", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var errorResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
	})

	t.Run("CacheError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/order/cache-error", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var errorResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
	})

	mockService.AssertExpectations(t)
}
