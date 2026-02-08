package services

import (
	"errors"
	"log"
	"smartfarm-api/config"
	"smartfarm-api/dto"
	"smartfarm-api/models"
	"smartfarm-api/repositories"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(req dto.CreateOrderRequest, userID uint) (dto.OrderResponse, error)
	GetMyOrders(userID uint) ([]dto.OrderResponse, error)
	GetAllOrders() ([]dto.OrderResponse, error) // For Admin/Farmer

	CreateSubscription(req dto.CreateSubscriptionRequest, userID uint) (dto.SubscriptionResponse, error)
	GetMySubscriptions(userID uint) ([]dto.SubscriptionResponse, error)
}

type orderService struct {
	orderRepo   repositories.OrderRepository
	productRepo repositories.ProductRepository
}

func NewOrderService(orderRepo repositories.OrderRepository, productRepo repositories.ProductRepository) OrderService {
	return &orderService{orderRepo, productRepo}
}

func (s *orderService) CreateOrder(req dto.CreateOrderRequest, userID uint) (dto.OrderResponse, error) {
	if len(req.Items) == 0 {
		return dto.OrderResponse{}, errors.New("cart is empty")
	}

	var createdOrder models.Order

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		txOrderRepo := s.orderRepo.WithTx(tx)
		txProductRepo := s.productRepo.WithTx(tx)

		var total float64
		var orderItems []models.OrderItem
		isPreOrder := false

		// 1. Process items and check stock
		for _, itemReq := range req.Items {
			log.Printf("[ORDER] User %d attempting to order Product %d, Qty %d", userID, itemReq.ProductID, itemReq.Quantity)

			// Find product with Lock (Select for Update) to prevent race conditions
			var product models.Product
			log.Printf("[LOCK] Acquiring lock for Product %d", itemReq.ProductID)
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, itemReq.ProductID).Error; err != nil {
				log.Printf("[ERROR] Product %d not found", itemReq.ProductID)
				return errors.New("product not found")
			}
			log.Printf("[LOCK] Lock acquired for Product %d, Current Stock: %d", itemReq.ProductID, product.Stock)

			if product.Stock < itemReq.Quantity {
				log.Printf("[REJECT] Insufficient stock for Product %d (Available: %d, Requested: %d)",
					itemReq.ProductID, product.Stock, itemReq.Quantity)
				return errors.New("insufficient stock for " + product.Name)
			}

			price := product.Price
			subTotal := price * float64(itemReq.Quantity)
			total += subTotal

			if product.IsPreOrder {
				isPreOrder = true
			}

			orderItems = append(orderItems, models.OrderItem{
				ProductID: product.ID,
				Quantity:  itemReq.Quantity,
				Price:     price,
			})

			// Update stock
			oldStock := product.Stock
			product.Stock -= itemReq.Quantity
			if err := txProductRepo.Update(&product); err != nil {
				log.Printf("[ERROR] Failed to update stock for Product %d", itemReq.ProductID)
				return err
			}
			log.Printf("[SUCCESS] Stock updated for Product %d (Old: %d, New: %d)",
				itemReq.ProductID, oldStock, product.Stock)
		}

		orderType := "regular"
		if isPreOrder {
			orderType = "preorder"
		}

		createdOrder = models.Order{
			UserID:     userID,
			TotalPrice: total,
			Status:     "pending",
			Type:       orderType,
			OrderItems: orderItems,
		}

		if req.AddressID != 0 {
			createdOrder.AddressID = &req.AddressID
		}

		// 2. Create Order
		if err := txOrderRepo.Create(&createdOrder); err != nil {
			log.Printf("[OrderService] DB Create Order Error: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		return dto.OrderResponse{}, err
	}

	return mapOrderToResponse(createdOrder), nil
}

func (s *orderService) GetMyOrders(userID uint) ([]dto.OrderResponse, error) {
	orders, err := s.orderRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.OrderResponse
	for _, o := range orders {
		responses = append(responses, mapOrderToResponse(o))
	}
	return responses, nil
}

func (s *orderService) GetAllOrders() ([]dto.OrderResponse, error) {
	orders, err := s.orderRepo.FindAll()
	if err != nil {
		return nil, err
	}
	var responses []dto.OrderResponse
	for _, o := range orders {
		responses = append(responses, mapOrderToResponse(o))
	}
	return responses, nil
}

func (s *orderService) CreateSubscription(req dto.CreateSubscriptionRequest, userID uint) (dto.SubscriptionResponse, error) {
	product, err := s.productRepo.FindByID(req.ProductID)
	if err != nil {
		return dto.SubscriptionResponse{}, errors.New("product not found")
	}

	if !product.IsSubscription {
		return dto.SubscriptionResponse{}, errors.New("this product is not available for subscription")
	}

	startDate := time.Now()
	var endDate time.Time

	// Simple Logic: Duration is number of weeks/months
	if req.Frequency == "weekly" {
		endDate = startDate.AddDate(0, 0, req.Duration*7)
	} else if req.Frequency == "monthly" {
		endDate = startDate.AddDate(0, req.Duration, 0)
	}

	sub := models.Subscription{
		UserID:    userID,
		ProductID: req.ProductID,
		Frequency: req.Frequency,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    "active",
	}

	err = s.orderRepo.CreateSubscription(&sub)
	if err != nil {
		return dto.SubscriptionResponse{}, err
	}

	// Re-fetch to get associations if needed or map manually
	sub.Product = product

	return mapSubscriptionToResponse(sub), nil
}

func (s *orderService) GetMySubscriptions(userID uint) ([]dto.SubscriptionResponse, error) {
	subs, err := s.orderRepo.FindSubscriptionsByUserID(userID)
	if err != nil {
		return nil, err
	}
	var responses []dto.SubscriptionResponse
	for _, sub := range subs {
		responses = append(responses, mapSubscriptionToResponse(sub))
	}
	return responses, nil
}

// Helpers
func mapOrderToResponse(o models.Order) dto.OrderResponse {
	var itemResponses []dto.OrderItemResponse
	for _, item := range o.OrderItems {
		productName := item.Product.Name
		if productName == "" && item.ProductID != 0 {
			// Fallback if not preloaded (though it should be)
			productName = "Product #" + strconv.FormatUint(uint64(item.ProductID), 10)
		}
		itemResponses = append(itemResponses, dto.OrderItemResponse{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Price:       item.Price,
			SubTotal:    item.Price * float64(item.Quantity),
		})
	}

	return dto.OrderResponse{
		ID:           o.ID,
		UserID:       o.UserID,
		TotalPrice:   o.TotalPrice,
		Status:       o.Status,
		Type:         o.Type,
		PaymentProof: o.PaymentProof,
		Items:        itemResponses,
		CreatedAt:    o.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func mapSubscriptionToResponse(s models.Subscription) dto.SubscriptionResponse {
	return dto.SubscriptionResponse{
		ID:          s.ID,
		UserID:      s.UserID,
		ProductName: s.Product.Name,
		Frequency:   s.Frequency,
		StartDate:   s.StartDate.Format("2006-01-02"),
		EndDate:     s.EndDate.Format("2006-01-02"),
		Status:      s.Status,
	}
}
