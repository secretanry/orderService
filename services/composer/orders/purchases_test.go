package orders

import (
	contextpkg "context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"wb-L0/modules/monitoring"
	"wb-L0/services/cache"
	"wb-L0/services/database"
	"wb-L0/structs"
)

// MockCache is a mock implementation of the cache service
type MockCache struct {
	mock.Mock
}

func (m *MockCache) GetOrder(ctx contextpkg.Context, orderId string) (*structs.Order, error) {
	args := m.Called(ctx, orderId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*structs.Order), args.Error(1)
}

func (m *MockCache) PutOrder(ctx contextpkg.Context, orderId string, order *structs.Order) error {
	args := m.Called(ctx, orderId, order)
	return args.Error(0)
}

func (m *MockCache) HealthCheck(ctx contextpkg.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockDatabase is a mock implementation of the database service
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetOrderById(ctx contextpkg.Context, orderId string) (*structs.Order, error) {
	args := m.Called(ctx, orderId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*structs.Order), args.Error(1)
}

func (m *MockDatabase) InsertOrder(ctx contextpkg.Context, order *structs.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockDatabase) HealthCheck(ctx contextpkg.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestGetOrderByIdCacheHit tests successful order retrieval from cache
func TestGetOrderByIdCacheHit(t *testing.T) {
	// Create mocks
	mockCache := new(MockCache)
	mockDB := new(MockDatabase)

	// Create test order
	expectedOrder := &structs.Order{
		OrderUid:    "cache-hit-test",
		TrackNumber: "CACHE_TRACK_123",
		CustomerId:  "test_customer",
	}

	// Set up mock expectations
	mockCache.On("GetOrder", mock.Anything, "cache-hit-test").Return(expectedOrder, nil)

	// Set up services
	cache.SetCache(mockCache)
	database.SetDatabase(mockDB)

	// Call function
	order, err := GetOrderById(contextpkg.Background(), "cache-hit-test")

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, expectedOrder.OrderUid, order.OrderUid)
	assert.Equal(t, expectedOrder.TrackNumber, order.TrackNumber)

	// Verify cache was called but database was not
	mockCache.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "GetOrderById")
}

// TestGetOrderByIdCacheMissDatabaseHit tests order retrieval when cache misses but database has the order
func TestGetOrderByIdCacheMissDatabaseHit(t *testing.T) {
	// Create mocks
	mockCache := new(MockCache)
	mockDB := new(MockDatabase)

	// Create test order
	expectedOrder := &structs.Order{
		OrderUid:    "cache-miss-test",
		TrackNumber: "DB_TRACK_456",
		CustomerId:  "test_customer",
	}

	// Set up mock expectations - the function calls GetCache().GetOrder() and GetDatabase().GetOrderById()
	mockCache.On("GetOrder", mock.Anything, "cache-miss-test").Return(nil, cache.ErrCacheMiss{Key: "cache-miss-test"})
	mockDB.On("GetOrderById", mock.Anything, "cache-miss-test").Return(expectedOrder, nil)
	mockCache.On("PutOrder", mock.Anything, "cache-miss-test", expectedOrder).Return(nil)

	// Set up services - this sets the global instances that the function uses
	cache.SetCache(mockCache)
	database.SetDatabase(mockDB)

	// Call function
	order, err := GetOrderById(contextpkg.Background(), "cache-miss-test")

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, expectedOrder.OrderUid, order.OrderUid)
	assert.Equal(t, expectedOrder.TrackNumber, order.TrackNumber)

	// Verify both cache and database were called
	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

// TestGetOrderByIdCacheMissDatabaseNotFound tests when cache misses and database doesn't have the order
func TestGetOrderByIdCacheMissDatabaseNotFound(t *testing.T) {
	// Create mocks
	mockCache := new(MockCache)
	mockDB := new(MockDatabase)

	// Set up mock expectations
	mockCache.On("GetOrder", mock.Anything, "not-found-test").Return(nil, cache.ErrCacheMiss{Key: "not-found-test"})
	mockDB.On("GetOrderById", mock.Anything, "not-found-test").Return(nil, assert.AnError)

	// Set up services
	cache.SetCache(mockCache)
	database.SetDatabase(mockDB)

	// Call function
	order, err := GetOrderById(contextpkg.Background(), "not-found-test")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, order)

	// Verify both cache and database were called
	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

// TestGetOrderByIdCacheError tests when cache returns an error (not cache miss)
func TestGetOrderByIdCacheError(t *testing.T) {
	// Create mocks
	mockCache := new(MockCache)
	mockDB := new(MockDatabase)

	// Set up mock expectations
	mockCache.On("GetOrder", mock.Anything, "cache-error-test").Return(nil, assert.AnError)

	// Set up services
	cache.SetCache(mockCache)
	database.SetDatabase(mockDB)

	// Call function
	order, err := GetOrderById(contextpkg.Background(), "cache-error-test")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, order)

	// Verify cache was called but database was not
	mockCache.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "GetOrderById")
}

// TestGetOrderByIdCachePutError tests when cache put operation fails
func TestGetOrderByIdCachePutError(t *testing.T) {
	// Create mocks
	mockCache := new(MockCache)
	mockDB := new(MockDatabase)

	// Create test order
	expectedOrder := &structs.Order{
		OrderUid:    "cache-put-error-test",
		TrackNumber: "DB_TRACK_789",
		CustomerId:  "test_customer",
	}

	// Set up mock expectations
	mockCache.On("GetOrder", mock.Anything, "cache-put-error-test").Return(nil, cache.ErrCacheMiss{Key: "cache-put-error-test"})
	mockDB.On("GetOrderById", mock.Anything, "cache-put-error-test").Return(expectedOrder, nil)
	mockCache.On("PutOrder", mock.Anything, "cache-put-error-test", expectedOrder).Return(assert.AnError)

	// Set up services
	cache.SetCache(mockCache)
	database.SetDatabase(mockDB)

	// Call function
	order, err := GetOrderById(contextpkg.Background(), "cache-put-error-test")

	// Assertions
	require.NoError(t, err) // Should still succeed even if cache put fails
	assert.Equal(t, expectedOrder.OrderUid, order.OrderUid)
	assert.Equal(t, expectedOrder.TrackNumber, order.TrackNumber)

	// Verify both cache and database were called
	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

// TestGetOrderByIdWithComplexOrder tests order retrieval with complex nested data
func TestGetOrderByIdWithComplexOrder(t *testing.T) {
	// Create mocks
	mockCache := new(MockCache)
	mockDB := new(MockDatabase)

	// Create complex test order
	expectedOrder := &structs.Order{
		OrderUid:          "complex-order-test",
		TrackNumber:       "COMPLEX_TRACK_999",
		Entry:             "WBILMCOMPLEX",
		Locale:            "ru",
		InternalSignature: "internal_sig",
		CustomerId:        "complex_customer",
		DeliveryService:   "russian_post",
		Shardkey:          "5",
		SmId:              123,
		DateCreated:       "2023-12-01T10:30:00Z",
		OofShard:          "2",
		Delivery: structs.Delivery{
			Name:    "Иван Иванов",
			Phone:   "+79001234567",
			Zip:     "123456",
			City:    "Москва",
			Address: "ул. Пушкина, д. 10",
			Region:  "Московская область",
			Email:   "ivan@example.com",
		},
		Payment: structs.Payment{
			Transaction:  "complex_transaction_123",
			RequestId:    "req_456",
			Currency:     "RUB",
			Provider:     "sberbank",
			Amount:       5000,
			PaymentDt:    1701437400,
			Bank:         "sberbank",
			DeliveryCost: 300,
			GoodsTotal:   4700,
			CustomFee:    100,
		},
		Items: []structs.Item{
			{
				ChrtId:      123456,
				TrackNumber: "COMPLEX_TRACK_999",
				Price:       2500,
				Rid:         "complex_rid_1",
				Name:        "Ноутбук",
				Sale:        0,
				Size:        "15.6",
				TotalPrice:  2500,
				NmId:        987654,
				Brand:       "Lenovo",
				Status:      200,
			},
			{
				ChrtId:      789012,
				TrackNumber: "COMPLEX_TRACK_999",
				Price:       2200,
				Rid:         "complex_rid_2",
				Name:        "Мышь",
				Sale:        10,
				Size:        "M",
				TotalPrice:  2200,
				NmId:        654321,
				Brand:       "Logitech",
				Status:      200,
			},
		},
	}

	// Set up mock expectations
	mockCache.On("GetOrder", mock.Anything, "complex-order-test").Return(nil, cache.ErrCacheMiss{Key: "complex-order-test"})
	mockDB.On("GetOrderById", mock.Anything, "complex-order-test").Return(expectedOrder, nil)
	mockCache.On("PutOrder", mock.Anything, "complex-order-test", expectedOrder).Return(nil)

	// Set up services
	cache.SetCache(mockCache)
	database.SetDatabase(mockDB)

	// Call function
	order, err := GetOrderById(contextpkg.Background(), "complex-order-test")

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, expectedOrder.OrderUid, order.OrderUid)
	assert.Equal(t, expectedOrder.Delivery.Name, order.Delivery.Name)
	assert.Equal(t, expectedOrder.Payment.Amount, order.Payment.Amount)
	assert.Len(t, order.Items, 2)
	assert.Equal(t, expectedOrder.Items[0].Name, order.Items[0].Name)
	assert.Equal(t, expectedOrder.Items[1].Brand, order.Items[1].Brand)

	// Verify both cache and database were called
	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

// TestGetOrderByIdContextCancellation tests behavior when context is cancelled
func TestGetOrderByIdContextCancellation(t *testing.T) {
	// Create mocks
	mockCache := new(MockCache)
	mockDB := new(MockDatabase)

	// Set up mock expectations for cancelled context
	mockCache.On("GetOrder", mock.Anything, "context-test").Return(nil, contextpkg.Canceled)

	// Set up services
	cache.SetCache(mockCache)
	database.SetDatabase(mockDB)

	// Create cancelled context
	ctx, cancel := contextpkg.WithCancel(contextpkg.Background())
	cancel() // Cancel immediately

	// Call function with cancelled context
	order, err := GetOrderById(ctx, "context-test")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, order)

	// Verify cache was called but database was not
	mockCache.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "GetOrderById")
}

// TestGetOrderByIdEmptyOrderId tests behavior with empty order ID
func TestGetOrderByIdEmptyOrderId(t *testing.T) {
	// Create mocks
	mockCache := new(MockCache)
	mockDB := new(MockDatabase)

	// Set up mock expectations for empty order ID
	mockCache.On("GetOrder", mock.Anything, "").Return(nil, cache.ErrCacheMiss{Key: ""})
	mockDB.On("GetOrderById", mock.Anything, "").Return(nil, assert.AnError)

	// Set up services
	cache.SetCache(mockCache)
	database.SetDatabase(mockDB)

	// Call function with empty order ID
	order, err := GetOrderById(contextpkg.Background(), "")

	// Assertion
	assert.Error(t, err)
	assert.Nil(t, order)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestMain(m *testing.M) {
	// Reset monitoring state before running tests
	monitoring.ResetForTesting()

	// Initialize monitoring in test mode
	monitoringInstance := &monitoring.Monitoring{}
	errChan := make(chan error, 1)
	if err := monitoringInstance.Init(errChan); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
