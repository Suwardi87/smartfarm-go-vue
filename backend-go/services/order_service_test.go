package services

import (
	"smartfarm-api/config"
	"smartfarm-api/dto"
	"smartfarm-api/models"
	"smartfarm-api/repositories"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConcurrentOrders_RaceCondition tests that concurrent orders don't cause overselling
func TestConcurrentOrders_RaceCondition(t *testing.T) {
	// Setup test database
	db := config.DB
	if db == nil {
		t.Skip("Database not available for testing")
	}

	// Clean up test data
	db.Exec("DELETE FROM order_items")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM products WHERE name LIKE 'Test Product%'")
	db.Exec("DELETE FROM users WHERE email LIKE 'test%@race.com'")

	// Create test product with stock = 1
	testProduct := models.Product{
		Name:        "Test Product - Race Condition",
		Description: "Product for testing race condition",
		Price:       100000,
		Stock:       1, // Only 1 item available
		Category:    "Test",
		FarmerID:    1, // Assuming farmer with ID 1 exists
		ImageURL:    "test.jpg",
	}
	db.Create(&testProduct)

	// Create 2 test users
	user1 := models.User{
		Name:     "Test User A",
		Email:    "testA@race.com",
		Password: "hashed",
		Role:     "pembeli",
	}
	user2 := models.User{
		Name:     "Test User B",
		Email:    "testB@race.com",
		Password: "hashed",
		Role:     "pembeli",
	}
	db.Create(&user1)
	db.Create(&user2)

	// Create test address
	address := models.Address{
		UserID:        user1.ID,
		Label:         "Test Address",
		RecipientName: "Test User",
		PhoneNumber:   "08123456789",
		Street:        "Test Street",
		City:          "Test City",
		Province:      "Test Province",
		PostalCode:    "12345",
		IsDefault:     true,
	}
	db.Create(&address)

	// Initialize services
	orderRepo := repositories.NewOrderRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderService := NewOrderService(orderRepo, productRepo)

	// Simulate 2 concurrent orders
	var wg sync.WaitGroup
	var successCount int
	var mu sync.Mutex
	var errors []error

	// Request payload
	createOrderRequest := dto.CreateOrderRequest{
		Items: []dto.OrderItemRequest{
			{
				ProductID: testProduct.ID,
				Quantity:  1,
			},
		},
		AddressID: address.ID,
	}

	// Goroutine 1 - User A orders
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := orderService.CreateOrder(createOrderRequest, user1.ID)
		mu.Lock()
		if err == nil {
			successCount++
		} else {
			errors = append(errors, err)
		}
		mu.Unlock()
	}()

	// Goroutine 2 - User B orders (almost simultaneously)
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := orderService.CreateOrder(createOrderRequest, user2.ID)
		mu.Lock()
		if err == nil {
			successCount++
		} else {
			errors = append(errors, err)
		}
		mu.Unlock()
	}()

	// Wait for both goroutines to complete
	wg.Wait()

	// ✅ ASSERTION 1: Only 1 order should succeed
	assert.Equal(t, 1, successCount, "Only 1 order should succeed when stock is 1")

	// ✅ ASSERTION 2: One error should occur
	assert.Equal(t, 1, len(errors), "One order should fail with error")

	// ✅ ASSERTION 3: Stock should be 0
	var updatedProduct models.Product
	db.First(&updatedProduct, testProduct.ID)
	assert.Equal(t, 0, updatedProduct.Stock, "Stock should be 0 after successful order")

	// ✅ ASSERTION 4: Only 1 order should exist in database
	var orderCount int64
	db.Model(&models.Order{}).Count(&orderCount)
	assert.Equal(t, int64(1), orderCount, "Only 1 order should exist in database")

	// Clean up
	db.Exec("DELETE FROM order_items")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM products WHERE id = ?", testProduct.ID)
	db.Exec("DELETE FROM users WHERE id IN (?, ?)", user1.ID, user2.ID)
	db.Exec("DELETE FROM addresses WHERE id = ?", address.ID)

	t.Log("✅ Race Condition Test PASSED - System correctly prevents overselling!")
}

// TestHighConcurrency_MultipleProducts tests system under high concurrent load
func TestHighConcurrency_MultipleProducts(t *testing.T) {
	db := config.DB
	if db == nil {
		t.Skip("Database not available for testing")
	}

	// Clean up
	db.Exec("DELETE FROM order_items")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM products WHERE name LIKE 'Concurrent Test%'")

	// Create product with stock = 10
	testProduct := models.Product{
		Name:        "Concurrent Test Product",
		Description: "Product for high concurrency testing",
		Price:       50000,
		Stock:       10, // 10 items available
		Category:    "Test",
		FarmerID:    1,
		ImageURL:    "test.jpg",
	}
	db.Create(&testProduct)

	// Initialize services
	orderRepo := repositories.NewOrderRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderService := NewOrderService(orderRepo, productRepo)

	// Simulate 20 concurrent orders (each ordering 1 item)
	// Expected: 10 succeed, 10 fail
	concurrentUsers := 20
	var wg sync.WaitGroup
	var successCount int
	var mu sync.Mutex

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userIndex int) {
			defer wg.Done()

			// Create temporary user
			user := models.User{
				Name:     "Concurrent User",
				Email:    "concurrent@test.com",
				Password: "hashed",
				Role:     "pembeli",
			}
			db.Create(&user)

			req := dto.CreateOrderRequest{
				Items: []dto.OrderItemRequest{
					{ProductID: testProduct.ID, Quantity: 1},
				},
			}

			_, err := orderService.CreateOrder(req, user.ID)
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}

			// Clean up user
			db.Delete(&user)
		}(i)
	}

	wg.Wait()

	// ✅ ASSERTION: Exactly 10 orders should succeed
	assert.Equal(t, 10, successCount, "Exactly 10 orders should succeed when stock is 10")

	// ✅ ASSERTION: Stock should be 0
	var updatedProduct models.Product
	db.First(&updatedProduct, testProduct.ID)
	assert.Equal(t, 0, updatedProduct.Stock, "Stock should be 0 after all successful orders")

	// Clean up
	db.Exec("DELETE FROM order_items")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM products WHERE id = ?", testProduct.ID)

	t.Log("✅ High Concurrency Test PASSED - System handles 20 concurrent requests correctly!")
}
