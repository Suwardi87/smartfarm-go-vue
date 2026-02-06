package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"smartfarm-api/config"
	"smartfarm-api/dto"
	"smartfarm-api/models"
	"smartfarm-api/repositories"
	"strings"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var paymentClient snap.Client
var orderRepo repositories.OrderRepository
var addressRepo repositories.AddressRepository

func InitPaymentService() {
	// Setup Midtrans
	paymentClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	// Initialize repositories
	orderRepo = repositories.NewOrderRepository(config.DB)
	addressRepo = repositories.NewAddressRepository(config.DB)
}

func CreatePayment(userID uint, req dto.CreatePaymentRequest) (*models.Payment, string, error) {
	log.Printf("[PaymentService] Starting CreatePayment for UserID: %d, OrderID: %d, Amount: %f", userID, req.OrderID, req.Amount)

	// Get order
	order, err := orderRepo.FindByID(req.OrderID)
	if err != nil {
		log.Printf("[PaymentService] Error: Order %d not found: %v", req.OrderID, err)
		return nil, "", errors.New("order not found")
	}
	log.Printf("[PaymentService] Order found: %+v", order)

	if order.UserID != userID {
		log.Printf("[PaymentService] Error: UserID %d unauthorized for Order %d (Owner: %d)", userID, req.OrderID, order.UserID)
		return nil, "", errors.New("unauthorized")
	}

	// Get address
	address, err := addressRepo.FindByID(req.AddressID)
	if err != nil {
		log.Printf("[PaymentService] Error: Address %d not found: %v", req.AddressID, err)
		return nil, "", errors.New("address not found: please add an address first")
	}
	log.Printf("[PaymentService] Address found: %+v", address)

	// Create payment record
	payment := models.Payment{
		OrderID:       req.OrderID,
		UserID:        userID,
		Amount:        req.Amount,
		Status:        "pending",
		TransactionID: fmt.Sprintf("ORD-%03d-%d", req.OrderID, time.Now().Unix()),
	}

	log.Printf("[PaymentService] Creating payment record in DB...")
	if err := repositories.CreatePayment(&payment); err != nil {
		log.Printf("[PaymentService] DB Error creating payment: %v", err)
		return nil, "", fmt.Errorf("failed to create payment record: %v", err)
	}
	log.Printf("[PaymentService] Payment record created with ID: %d", payment.ID)

	// Midtrans Logic
	serverKey := strings.TrimSpace(os.Getenv("MIDTRANS_SERVER_KEY"))
	var snapToken string
	var redirectURL string

	// 🕵️ Check for Placeholder Key (MOCK MODE)
	if serverKey == "" || serverKey == "YOUR_MIDTRANS_SERVER_KEY_HERE" || strings.Contains(serverKey, "PLACEHOLDER") {
		log.Println("🛠️ [PaymentService] Using MOCK MODE (Key is missing or placeholder)")
		snapToken = fmt.Sprintf("mock-token-%d-%d", payment.ID, time.Now().Unix())
		redirectURL = "https://app.sandbox.midtrans.com/snap/v2/vtweb/" + snapToken
	} else {
		log.Printf("[PaymentService] Using REAL MIDTRANS MODE (Key starts with %s...)", serverKey[:4])
		// Real Midtrans transaction
		items := []midtrans.ItemDetails{
			{
				ID:    fmt.Sprintf("ORD-%d", req.OrderID),
				Price: int64(req.Amount),
				Qty:   1,
				Name:  fmt.Sprintf("Produk SmartFarm #%d", req.OrderID),
			},
		}
		snapReq := snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  payment.TransactionID,
				GrossAmt: int64(req.Amount),
			},
			CustomerDetail: &midtrans.CustomerDetails{
				FName: address.RecipientName,
				Phone: address.PhoneNumber,
			},
			Items: &items,
		}

		log.Printf("[PaymentService] Calling Midtrans API...")
		snapResp, err := paymentClient.CreateTransaction(&snapReq)
		if err != nil {
			log.Printf("[PaymentService] Midtrans API Error: %v", err)
			return nil, "", fmt.Errorf("midtrans error: %v (check your server key)", err)
		}
		log.Printf("[PaymentService] Midtrans Response received. Token: %s", snapResp.Token)
		snapToken = snapResp.Token
		redirectURL = snapResp.RedirectURL
	}

	// Update payment with snap token
	payment.SnapToken = snapToken
	payment.SnapURL = redirectURL

	log.Printf("[PaymentService] Updating payment record in DB with token...")
	if err := repositories.UpdatePayment(&payment); err != nil {
		log.Printf("[PaymentService] DB Error updating payment: %v", err)
		return nil, "", fmt.Errorf("failed to update payment with token: %v", err)
	}

	// Update order with payment ID and address ID
	log.Printf("[PaymentService] Linking payment %d to order %d in DB...", payment.ID, order.ID)
	if err := orderRepo.UpdatePaymentInfo(order.ID, payment.ID, req.AddressID); err != nil {
		log.Printf("[PaymentService] DB Error linking order to payment: %v", err)
		return nil, "", fmt.Errorf("failed to link order to payment: %v", err)
	}

	log.Printf("[PaymentService] Successfully completed CreatePayment for ID: %d", payment.ID)
	return &payment, snapToken, nil
}

func ProcessPaymentWebhook(webhookData map[string]interface{}) error {
	transactionID := webhookData["order_id"].(string)
	transactionStatus := webhookData["transaction_status"].(string)

	payment, err := repositories.GetPaymentByTransactionID(transactionID)
	if err != nil {
		return errors.New("payment not found")
	}

	// Update payment status
	switch transactionStatus {
	case "capture", "settlement":
		payment.Status = "success"
	case "pending":
		payment.Status = "pending"
	case "deny", "cancel", "expire":
		payment.Status = "failed"
	}

	if err := repositories.UpdatePayment(payment); err != nil {
		return err
	}

	// Update order status
	order, err := orderRepo.FindByID(payment.OrderID)
	if err != nil {
		return err
	}

	if payment.Status == "success" {
		order.Status = "paid"
	} else if payment.Status == "failed" {
		order.Status = "cancelled"
	}

	return orderRepo.Update(&order)
}

func GetPaymentByOrderID(userID uint, orderID uint) (*models.Payment, error) {
	payment, err := repositories.GetPaymentByOrderID(orderID)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	// Verify ownership
	if payment.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return payment, nil
}

func ConfirmMockPayment(paymentID uint) error {
	serverKey := strings.TrimSpace(os.Getenv("MIDTRANS_SERVER_KEY"))
	if serverKey != "" && serverKey != "YOUR_MIDTRANS_SERVER_KEY_HERE" && !strings.Contains(serverKey, "PLACEHOLDER") {
		return errors.New("cannot confirm mock payment in production/real mode")
	}

	payment, err := repositories.GetPaymentByID(paymentID)
	if err != nil {
		return errors.New("payment not found")
	}

	payment.Status = "success"
	if err := repositories.UpdatePayment(payment); err != nil {
		return err
	}

	order, err := orderRepo.FindByID(payment.OrderID)
	if err != nil {
		return err
	}

	order.Status = "paid"
	return orderRepo.Update(&order)
}
