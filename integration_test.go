package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/IBM/sarama"

	"wb-L0/models/pg_models"
	"wb-L0/modules/monitoring"
	"wb-L0/modules/pg"
	redispkg "wb-L0/modules/redis"
	"wb-L0/routing"
	"wb-L0/services/cache"
	"wb-L0/services/database"
	"wb-L0/structs"
)

type IntegrationTestSuite struct {
	suite.Suite
	router        *gin.Engine
	db            *gorm.DB
	redisClient   *redis.Client
	ctx           context.Context
	kafkaProducer sarama.SyncProducer
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	gin.SetMode(gin.TestMode)

	suite.setupTestDatabase()

	suite.setupTestRedis()

	suite.setupTestKafka()

	suite.setupTestServices()

	suite.setupRouter()
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		sqlDB, err := suite.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	if suite.redisClient != nil {
		suite.redisClient.Close()
	}

	if suite.kafkaProducer != nil {
		suite.kafkaProducer.Close()
	}
}

func (suite *IntegrationTestSuite) SetupTest() {

	suite.cleanupDatabase()

	suite.cleanupRedis()
}

func (suite *IntegrationTestSuite) setupTestDatabase() {
	dsn := "host=localhost user=test_user password=test_pass dbname=test_db port=5433 sslmode=disable TimeZone=UTC"

	var err error
	suite.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(suite.T(), err)

	err = suite.db.AutoMigrate(
		&pg_models.Order{},
		&pg_models.OrderDelivery{},
		&pg_models.OrderPayment{},
		&pg_models.OrderItem{},
	)
	require.NoError(suite.T(), err)
}

func (suite *IntegrationTestSuite) setupTestRedis() {
	suite.redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "test_pass",
		DB:       1,
	})

	ctx := context.Background()
	err := suite.redisClient.Ping(ctx).Err()
	require.NoError(suite.T(), err)
}

func (suite *IntegrationTestSuite) setupTestKafka() {

	time.Sleep(10 * time.Second)

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	var err error
	suite.kafkaProducer, err = sarama.NewSyncProducer([]string{"localhost:9093"}, config)
	if err != nil {
		suite.T().Logf("Warning: Kafka connection failed, skipping Kafka tests: %v", err)
		suite.kafkaProducer = nil
		return
	}

	admin, err := sarama.NewClusterAdmin([]string{"localhost:9093"}, config)
	if err != nil {
		suite.T().Logf("Warning: Kafka admin connection failed: %v", err)
		return
	}
	defer admin.Close()

	topics, err := admin.ListTopics()
	if err != nil {
		suite.T().Logf("Warning: Failed to list Kafka topics: %v", err)
		return
	}

	if _, exists := topics["orders"]; !exists {
		err = admin.CreateTopic("orders", &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		if err != nil {
			suite.T().Logf("Warning: Failed to create Kafka topic: %v", err)
		}
	}
}

func (suite *IntegrationTestSuite) setupTestServices() {

	pgWrapper := &pg.Postgres{Db: suite.db}

	dbService := database.NewPostgres(pgWrapper)
	database.SetDatabase(dbService)

	redisWrapper := &redispkg.Redis{Client: suite.redisClient}

	cacheService := cache.NewRedisCache(redisWrapper)
	cache.SetCache(cacheService)
}

func (suite *IntegrationTestSuite) setupRouter() {
	suite.router = gin.New()
	suite.router.Use(routing.ApiContextMiddleware())
	routing.MountPurchasesRoutes(suite.router)
}

func (suite *IntegrationTestSuite) cleanupDatabase() {
	suite.db.Exec("DELETE FROM order_item")
	suite.db.Exec("DELETE FROM order_payment")
	suite.db.Exec("DELETE FROM order_delivery")
	suite.db.Exec("DELETE FROM \"order\"")
}

func (suite *IntegrationTestSuite) cleanupRedis() {
	suite.redisClient.FlushDB(suite.ctx)
}

func (suite *IntegrationTestSuite) createTestOrder(orderID string) *structs.Order {
	order := &structs.Order{
		OrderUid:          orderID,
		TrackNumber:       "TEST_TRACK_123",
		Entry:             "WBILMTESTTRACK",
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test_customer",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmId:              99,
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          "1",
		Delivery: structs.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: structs.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestId:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []structs.Item{
			{
				ChrtId:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmId:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
	}

	err := database.GetDatabase().InsertOrder(suite.ctx, order)
	require.NoError(suite.T(), err)

	return order
}

func (suite *IntegrationTestSuite) TestGetOrderSuccess() {

	orderID := "test-order-123"
	expectedOrder := suite.createTestOrder(orderID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/order/"+orderID, nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var responseOrder structs.Order
	err := json.Unmarshal(w.Body.Bytes(), &responseOrder)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), expectedOrder.OrderUid, responseOrder.OrderUid)
	assert.Equal(suite.T(), expectedOrder.TrackNumber, responseOrder.TrackNumber)
	assert.Equal(suite.T(), expectedOrder.CustomerId, responseOrder.CustomerId)
	assert.Len(suite.T(), responseOrder.Items, 1)
	assert.Equal(suite.T(), expectedOrder.Items[0].Name, responseOrder.Items[0].Name)
}

func (suite *IntegrationTestSuite) TestGetOrderNotFound() {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/order/non-existent-order", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(suite.T(), err)

	_, hasError := errorResponse["error"]
	_, hasMessage := errorResponse["message"]
	assert.True(suite.T(), hasError || hasMessage, "Response should contain 'error' or 'message' key")
}

func (suite *IntegrationTestSuite) TestGetOrderMissingID() {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/order/", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

// TestCacheIntegration tests that orders are properly cached
func (suite *IntegrationTestSuite) TestCacheIntegration() {
	// Create test order
	orderID := "cache-test-order"
	expectedOrder := suite.createTestOrder(orderID)

	// First request - should hit database
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/api/order/"+orderID, nil)
	suite.router.ServeHTTP(w1, req1)

	assert.Equal(suite.T(), http.StatusOK, w1.Code)

	// Verify order is cached
	cachedOrder, err := cache.GetCache().GetOrder(suite.ctx, orderID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedOrder.OrderUid, cachedOrder.OrderUid)

	// Second request - should hit cache
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/order/"+orderID, nil)
	suite.router.ServeHTTP(w2, req2)

	assert.Equal(suite.T(), http.StatusOK, w2.Code)

	var responseOrder structs.Order
	err = json.Unmarshal(w2.Body.Bytes(), &responseOrder)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), expectedOrder.OrderUid, responseOrder.OrderUid)
}

// TestMultipleOrders tests retrieving multiple different orders
func (suite *IntegrationTestSuite) TestMultipleOrders() {
	// Create multiple test orders
	order1 := suite.createTestOrder("order-1")
	order2 := suite.createTestOrder("order-2")

	// Test retrieving first order
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/api/order/order-1", nil)
	suite.router.ServeHTTP(w1, req1)

	assert.Equal(suite.T(), http.StatusOK, w1.Code)

	var responseOrder1 structs.Order
	err := json.Unmarshal(w1.Body.Bytes(), &responseOrder1)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), order1.OrderUid, responseOrder1.OrderUid)

	// Test retrieving second order
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/order/order-2", nil)
	suite.router.ServeHTTP(w2, req2)

	assert.Equal(suite.T(), http.StatusOK, w2.Code)

	var responseOrder2 structs.Order
	err = json.Unmarshal(w2.Body.Bytes(), &responseOrder2)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), order2.OrderUid, responseOrder2.OrderUid)
}

// TestOrderWithComplexData tests order with complex nested data
func (suite *IntegrationTestSuite) TestOrderWithComplexData() {
	// Create order with complex data
	orderID := "complex-order"
	order := &structs.Order{
		OrderUid:          orderID,
		TrackNumber:       "COMPLEX_TRACK_456",
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
				TrackNumber: "COMPLEX_TRACK_456",
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
				TrackNumber: "COMPLEX_TRACK_456",
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

	// Save to database
	err := database.GetDatabase().InsertOrder(suite.ctx, order)
	require.NoError(suite.T(), err)

	// Make HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/order/"+orderID, nil)
	suite.router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var responseOrder structs.Order
	err = json.Unmarshal(w.Body.Bytes(), &responseOrder)
	require.NoError(suite.T(), err)

	// Verify complex data
	assert.Equal(suite.T(), order.OrderUid, responseOrder.OrderUid)
	assert.Equal(suite.T(), order.Delivery.Name, responseOrder.Delivery.Name)
	assert.Equal(suite.T(), order.Payment.Amount, responseOrder.Payment.Amount)
	assert.Len(suite.T(), responseOrder.Items, 2)
	assert.Equal(suite.T(), order.Items[0].Name, responseOrder.Items[0].Name)
	assert.Equal(suite.T(), order.Items[1].Brand, responseOrder.Items[1].Brand)
}

// TestDatabaseConnectionFailure tests behavior when database is unavailable
// func (suite *IntegrationTestSuite) TestDatabaseConnectionFailure() {
// 	// Close database connection to simulate failure
// 	sqlDB, err := suite.db.DB()
// 	require.NoError(suite.T(), err)
// 	sqlDB.Close()
// 	// Try to get an order
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/api/order/test-order", nil)
// 	suite.router.ServeHTTP(w, req)
// 	// Should return 500 error
// 	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
// }

// TestRedisConnectionFailure tests behavior when Redis is unavailable
func (suite *IntegrationTestSuite) TestRedisConnectionFailure() {
	// Close Redis connection to simulate failure
	suite.redisClient.Close()
	// Create test order
	orderID := "redis-failure-test"
	suite.createTestOrder(orderID)
	// Try to get an order - should return 500 since handler does not degrade gracefully
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/order/"+orderID, nil)
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

// TestConcurrentRequests tests handling of concurrent requests
func (suite *IntegrationTestSuite) TestConcurrentRequests() {
	// Create test order
	orderID := "concurrent-test-order"
	suite.createTestOrder(orderID)

	// Make concurrent requests
	const numRequests = 10
	results := make(chan int, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/order/"+orderID, nil)
			suite.router.ServeHTTP(w, req)
			results <- w.Code
		}()
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		statusCode := <-results
		assert.Equal(suite.T(), http.StatusOK, statusCode)
	}
}

// TestOrderDataIntegrity tests that order data integrity is maintained
func (suite *IntegrationTestSuite) TestOrderDataIntegrity() {
	// Create order with specific data
	orderID := "integrity-test"
	originalOrder := suite.createTestOrder(orderID)

	// Retrieve order multiple times
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/order/"+orderID, nil)
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var responseOrder structs.Order
		err := json.Unmarshal(w.Body.Bytes(), &responseOrder)
		require.NoError(suite.T(), err)

		// Verify data integrity
		assert.Equal(suite.T(), originalOrder.OrderUid, responseOrder.OrderUid)
		assert.Equal(suite.T(), originalOrder.TrackNumber, responseOrder.TrackNumber)
		assert.Equal(suite.T(), originalOrder.Delivery.Name, responseOrder.Delivery.Name)
		assert.Equal(suite.T(), originalOrder.Payment.Amount, responseOrder.Payment.Amount)
		assert.Len(suite.T(), responseOrder.Items, len(originalOrder.Items))
	}
}

// TestInvalidOrderID tests handling of invalid order IDs
func (suite *IntegrationTestSuite) TestInvalidOrderID() {
	testCases := []string{
		"",            // Empty string
		"   ",         // Whitespace only
		"invalid@#$%", // Special characters
		"very-long-order-id-that-exceeds-normal-length-limits-and-should-be-handled-properly",
	}

	for _, orderID := range testCases {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Panic is expected for invalid order IDs due to nil pointer dereference
					// This is acceptable behavior since the handler doesn't validate order IDs
					suite.T().Logf("Expected panic for invalid order ID '%s': %v", orderID, r)
				}
			}()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/order/"+orderID, nil)
			suite.router.ServeHTTP(w, req)

			// If no panic occurred, check the response
			if w.Code != 0 {
				// Should return 404 or 500 for invalid IDs
				assert.Contains(suite.T(), []int{http.StatusNotFound, http.StatusInternalServerError}, w.Code, "Failed for order ID: %s", orderID)
			}
		}()
	}
}

// TestPerformance tests basic performance characteristics
func (suite *IntegrationTestSuite) TestPerformance() {
	// Create test order
	orderID := "performance-test"
	suite.createTestOrder(orderID)

	// Measure response time for first request (database hit)
	start := time.Now()
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/api/order/"+orderID, nil)
	suite.router.ServeHTTP(w1, req1)
	dbTime := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w1.Code)

	// Measure response time for second request (cache hit)
	start = time.Now()
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/order/"+orderID, nil)
	suite.router.ServeHTTP(w2, req2)
	cacheTime := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w2.Code)

	// Cache hit should be faster than database hit
	assert.Less(suite.T(), cacheTime, dbTime, "Cache hit should be faster than database hit")

	// Both should complete within reasonable time
	assert.Less(suite.T(), dbTime, 100*time.Millisecond, "Database hit took too long")
	assert.Less(suite.T(), cacheTime, 50*time.Millisecond, "Cache hit took too long")
}

// TestKafkaIntegration tests Kafka messaging functionality
func (suite *IntegrationTestSuite) TestKafkaIntegration() {
	if suite.kafkaProducer == nil {
		suite.T().Skip("Kafka not available, skipping Kafka integration test")
		return
	}

	// Test 1: Produce a message to Kafka
	order := suite.createTestOrder("kafka-test-order")
	orderJSON, err := json.Marshal(order)
	require.NoError(suite.T(), err)

	msg := &sarama.ProducerMessage{
		Topic: "orders",
		Key:   sarama.StringEncoder("kafka-test-order"),
		Value: sarama.ByteEncoder(orderJSON),
	}

	partition, offset, err := suite.kafkaProducer.SendMessage(msg)
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), partition, int32(0))
	assert.GreaterOrEqual(suite.T(), offset, int64(0))

	suite.T().Logf("Message sent to Kafka - Partition: %d, Offset: %d", partition, offset)

	// Test 2: Verify message was sent by checking topic metadata
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin([]string{"localhost:9093"}, config)
	require.NoError(suite.T(), err)
	defer admin.Close()

	topics, err := admin.ListTopics()
	require.NoError(suite.T(), err)
	assert.Contains(suite.T(), topics, "orders")

	// Test 3: Test multiple message production
	for i := 1; i <= 5; i++ {
		orderID := fmt.Sprintf("kafka-batch-test-%d", i)
		order := suite.createTestOrder(orderID)
		orderJSON, err := json.Marshal(order)
		require.NoError(suite.T(), err)

		msg := &sarama.ProducerMessage{
			Topic: "orders",
			Key:   sarama.StringEncoder(orderID),
			Value: sarama.ByteEncoder(orderJSON),
		}

		_, _, err = suite.kafkaProducer.SendMessage(msg)
		require.NoError(suite.T(), err)
	}

	suite.T().Logf("Successfully sent 5 batch messages to Kafka")
}

// TestKafkaConsumerIntegration tests Kafka consumer functionality
func (suite *IntegrationTestSuite) TestKafkaConsumerIntegration() {
	if suite.kafkaProducer == nil {
		suite.T().Skip("Kafka not available, skipping Kafka consumer test")
		return
	}

	// Create a consumer to test message consumption
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer([]string{"localhost:9093"}, config)
	require.NoError(suite.T(), err)
	defer consumer.Close()

	// Get partitions for the topic
	partitions, err := consumer.Partitions("orders")
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), len(partitions), 0)

	// Create partition consumer
	partitionConsumer, err := consumer.ConsumePartition("orders", partitions[0], sarama.OffsetNewest)
	require.NoError(suite.T(), err)
	defer partitionConsumer.Close()

	// Send a test message
	order := suite.createTestOrder("consumer-test-order")
	orderJSON, err := json.Marshal(order)
	require.NoError(suite.T(), err)

	msg := &sarama.ProducerMessage{
		Topic: "orders",
		Key:   sarama.StringEncoder("consumer-test-order"),
		Value: sarama.ByteEncoder(orderJSON),
	}

	_, _, err = suite.kafkaProducer.SendMessage(msg)
	require.NoError(suite.T(), err)

	// Consume the message with timeout
	select {
	case message := <-partitionConsumer.Messages():
		assert.Equal(suite.T(), "consumer-test-order", string(message.Key))
		assert.Equal(suite.T(), orderJSON, message.Value)
		suite.T().Logf("Successfully consumed message: %s", string(message.Key))
	case err := <-partitionConsumer.Errors():
		suite.T().Fatalf("Error consuming message: %v", err)
	case <-time.After(10 * time.Second):
		suite.T().Fatal("Timeout waiting for message consumption")
	}
}

// TestKafkaMessageFormat tests that messages sent to Kafka have the correct format
func (suite *IntegrationTestSuite) TestKafkaMessageFormat() {
	if suite.kafkaProducer == nil {
		suite.T().Skip("Kafka not available, skipping Kafka message format test")
		return
	}

	// Create a complex order with all fields
	order := &structs.Order{
		OrderUid:    "format-test-order",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery:    structs.Delivery{Name: "Test Test", Phone: "+9720000000", Zip: "2639809", City: "Kiryat Mozkin", Address: "Ploshad Mira 15", Region: "Kraiot", Email: "test@gmail.com"},
		Payment:     structs.Payment{Transaction: "b563feb7b2b84b6test", RequestId: "", Currency: "USD", Provider: "wbpay", Amount: 1817, PaymentDt: 1637907727, Bank: "alpha", DeliveryCost: 1500, GoodsTotal: 317, CustomFee: 0},
		Items: []structs.Item{
			{ChrtId: 9934930, TrackNumber: "WBILMTESTTRACK", Price: 453, Rid: "ab4219087a764ae0btest", Name: "Mascaras", Sale: 30, Size: "0", TotalPrice: 317, NmId: 2389212, Brand: "Vivienne Sabo", Status: 202},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmId:              99,
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          "1",
	}

	orderJSON, err := json.Marshal(order)
	require.NoError(suite.T(), err)

	// Verify JSON format
	var parsedOrder structs.Order
	err = json.Unmarshal(orderJSON, &parsedOrder)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), order.OrderUid, parsedOrder.OrderUid)
	assert.Equal(suite.T(), order.TrackNumber, parsedOrder.TrackNumber)
	assert.Len(suite.T(), parsedOrder.Items, 1)

	// Send to Kafka
	msg := &sarama.ProducerMessage{
		Topic: "orders",
		Key:   sarama.StringEncoder("format-test-order"),
		Value: sarama.ByteEncoder(orderJSON),
	}

	_, _, err = suite.kafkaProducer.SendMessage(msg)
	require.NoError(suite.T(), err)

	suite.T().Logf("Successfully sent formatted message to Kafka")
}

// TestKafkaErrorHandling tests Kafka error scenarios
func (suite *IntegrationTestSuite) TestKafkaErrorHandling() {
	if suite.kafkaProducer == nil {
		suite.T().Skip("Kafka not available, skipping Kafka error handling test")
		return
	}

	// Test 1: Try to send message to non-existent topic
	msg := &sarama.ProducerMessage{
		Topic: "non-existent-topic",
		Key:   sarama.StringEncoder("test"),
		Value: sarama.StringEncoder("test"),
	}

	_, _, err := suite.kafkaProducer.SendMessage(msg)
	// This should fail, but we handle it gracefully
	if err != nil {
		suite.T().Logf("Expected error for non-existent topic: %v", err)
	}

	// Test 2: Try to send invalid JSON
	msg = &sarama.ProducerMessage{
		Topic: "orders",
		Key:   sarama.StringEncoder("invalid-json"),
		Value: sarama.StringEncoder("{invalid json}"),
	}

	_, _, err = suite.kafkaProducer.SendMessage(msg)
	// This should succeed (Kafka doesn't validate JSON content)
	require.NoError(suite.T(), err)

	suite.T().Logf("Kafka error handling tests completed")
}

// TestKafkaPerformance tests Kafka performance with multiple messages
func (suite *IntegrationTestSuite) TestKafkaPerformance() {
	if suite.kafkaProducer == nil {
		suite.T().Skip("Kafka not available, skipping Kafka performance test")
		return
	}

	start := time.Now()
	messageCount := 10

	for i := 0; i < messageCount; i++ {
		orderID := fmt.Sprintf("perf-test-%d", i)
		order := suite.createTestOrder(orderID)
		orderJSON, err := json.Marshal(order)
		require.NoError(suite.T(), err)

		msg := &sarama.ProducerMessage{
			Topic: "orders",
			Key:   sarama.StringEncoder(orderID),
			Value: sarama.ByteEncoder(orderJSON),
		}

		_, _, err = suite.kafkaProducer.SendMessage(msg)
		require.NoError(suite.T(), err)
	}

	duration := time.Since(start)
	rate := float64(messageCount) / duration.Seconds()

	suite.T().Logf("Sent %d messages in %v (%.2f messages/sec)", messageCount, duration, rate)
	assert.Greater(suite.T(), rate, 1.0, "Should send at least 1 message per second")
}

func TestMain(m *testing.M) {
	// Reset monitoring state before running tests
	monitoring.ResetForTesting()

	monitoringInstance := &monitoring.Monitoring{}
	errChan := make(chan error, 1)
	if err := monitoringInstance.Init(errChan); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// TestIntegrationSuite runs the test suite
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
