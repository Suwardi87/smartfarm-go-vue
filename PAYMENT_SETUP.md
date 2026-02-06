# SmartFarm Payment Gateway Integration - Setup Guide

## Overview
This project integrates **Midtrans** as the payment gateway for secure online payments. Midtrans is Indonesia's leading payment aggregator supporting multiple payment methods (credit cards, debit cards, bank transfers, e-wallets, etc.).

## Architecture

### Backend (Go)
- **Payment Service**: Handles payment creation and webhook processing
- **Payment Repository**: Database operations for payment records
- **Payment Controller**: HTTP endpoints for payment operations
- **Models**: Payment and Order models with relationships

### Frontend (Vue 3)
- **Payment Service**: Midtrans Snap.js integration
- **Checkout Page**: Order review, address selection, payment initiation
- **Routes**: Protected checkout route

## Setup Instructions

### 1. Backend Configuration

#### 1.1 Get Midtrans Credentials
1. Go to https://dashboard.midtrans.com
2. Create a merchant account (or use existing one)
3. Go to Settings → Access Keys
4. Copy your **Server Key** and **Client Key**
5. Note the **Environment**: Sandbox (for testing) or Production

#### 1.2 Set Environment Variables
Edit `smartfarm-api/.env`:
```env
MIDTRANS_SERVER_KEY=your_server_key_here
MIDTRANS_CLIENT_KEY=your_client_key_here
```

#### 1.3 Database Migration
The Order model already has the following new fields:
- `PaymentID` - foreign key to Payment
- `AddressID` - foreign key to Address

Run your database migrations (GORM will auto-create the Payment table):
```bash
cd smartfarm-api
go run cmd/main.go
```

### 2. Frontend Configuration

#### 2.1 Set Midtrans Client Key
Edit `vue-tailwind-admin-dashboard-main/.env`:
```env
VITE_MIDTRANS_CLIENT_KEY=your_client_key_here
```

#### 2.2 Install Dependencies
The project uses axios which is already configured:
```bash
cd vue-tailwind-admin-dashboard-main
npm install
```

## API Endpoints

### Create Payment
**POST** `/payments` (Protected)

Request Body:
```json
{
  "order_id": 1,
  "address_id": 5,
  "amount": 150000
}
```

Response:
```json
{
  "data": {
    "payment_id": 1,
    "snap_token": "0ee2a02d-30ed-41f3-96c1-ccc9d806e7cd",
    "amount": 150000
  }
}
```

### Get Payment Status
**GET** `/payments/orders/:order_id` (Protected)

Response:
```json
{
  "data": {
    "payment_id": 1,
    "status": "success|pending|failed|expired",
    "amount": 150000
  }
}
```

### Payment Webhook
**POST** `/payments/webhook` (Public)

Midtrans will send POST requests to this endpoint to notify payment status changes.

## Frontend Flow

### Checkout Page (`/checkout`)

1. **Display Order Summary**
   - List all items in cart with quantities and prices
   - Calculate subtotal, tax, and total

2. **Select Delivery Address**
   - Show user's saved addresses
   - Auto-select default address
   - Option to manage addresses

3. **Initiate Payment**
   - User clicks "Pay Now" button
   - Frontend calls `POST /payments` with order details
   - Backend returns Midtrans snap token

4. **Midtrans Snap Payment Modal**
   - Frontend loads Midtrans Snap.js library
   - Display payment modal with token
   - User selects payment method
   - User completes payment

5. **Payment Status Handling**
   - **Success**: Clear cart, redirect to Orders page
   - **Pending**: Show pending message, user can check status later
   - **Failed**: Show error message, allow retry

## Payment Status Flow

```
Order Created (status: pending)
    ↓
Payment Created (status: pending)
    ↓
User initiates payment → Midtrans Snap
    ↓
Payment completed in Midtrans
    ↓
Webhook received
    ↓
Update Payment Status
    ↓
Update Order Status
    ↓
Success: order status = "paid"
```

## Payment Statuses

### Payment Model Statuses
- `pending` - Payment created, waiting for user action
- `success` - Payment verified by Midtrans
- `failed` - Payment rejected or cancelled
- `expired` - Payment window expired

### Order Model Status Updates
- When Payment Status = "success" → Order Status = "paid"
- When Payment Status = "failed" → Order Status = "cancelled"

## Testing Payment

### Using Midtrans Sandbox Credentials

1. **Test Card Numbers** (Sandbox only):
   - Visa Success: `4111 1111 1111 1111`
   - MasterCard Success: `5555 5555 5555 4444`
   - 3D Secure: Will prompt for OTP (use any 6 digits)

2. **Test Flow**:
   1. Add products to cart on homepage
   2. Go to Cart page
   3. Click "Proceed to Checkout"
   4. Select/create a delivery address
   5. Click "Pay Now"
   6. Choose payment method
   7. Enter test card details
   8. Complete payment

3. **Verify Payment**:
   - Check Orders page - status should show as "paid"
   - Check Midtrans Dashboard for transaction record

## Webhook Configuration

### Enable Midtrans Webhook
1. Go to https://dashboard.midtrans.com
2. Settings → Configuration
3. Set Notification URL: `https://yourdomain.com/payments/webhook`
4. Save

### Webhook Payload Handling
The webhook endpoint (`/payments/webhook`) receives Midtrans transaction updates and:
- Updates Payment record status
- Updates Order record status accordingly
- No response needed (endpoint returns 200 OK)

## Security Considerations

1. **Server Key**: Keep MIDTRANS_SERVER_KEY private (backend only)
2. **Client Key**: Can be exposed (frontend only)
3. **Signature Verification**: Production should verify Midtrans signatures on webhooks
4. **Address Validation**: Always verify address ownership before payment
5. **Order Verification**: Verify order belongs to user making payment

## Production Deployment

1. **Switch to Production**
   - Update payment_service.go: Change `midtrans.Sandbox` to `midtrans.Production`
   - Update environment variables with production keys

2. **HTTPS Required**
   - Midtrans requires HTTPS in production
   - Update webhook URL to HTTPS

3. **Database Backup**
   - Ensure Payment table is backed up regularly

## Troubleshooting

### Issue: "Failed to load Snap script"
- **Cause**: Network issue or Client Key is invalid
- **Solution**: Check Client Key in .env, verify internet connection

### Issue: "order not found"
- **Cause**: Order ID doesn't exist or user doesn't own it
- **Solution**: Ensure order was created before payment

### Issue: "address not found"
- **Cause**: Address ID is invalid
- **Solution**: User must create address before checkout

### Issue: Webhook not received
- **Cause**: Webhook URL not configured or DNS not pointing correctly
- **Solution**: Check Midtrans dashboard webhook configuration, verify domain

## File Structure

```
Backend:
smartfarm-api/
├── models/payment.go              # Payment model
├── repositories/payment_repository.go     # Payment CRUD
├── services/payment_service.go    # Midtrans integration
├── controllers/payment_controller.go     # HTTP endpoints
├── dto/payment.go                 # Request/Response DTOs
├── routes/routes.go               # Route definitions
└── .env                          # Configuration

Frontend:
vue-tailwind-admin-dashboard-main/
├── src/
│   ├── services/paymentService.ts        # Snap.js integration
│   ├── dto/payment/index.ts               # Type definitions
│   ├── views/Marketplace/Checkout.vue    # Checkout page
│   └── .env                               # Configuration
└── index.html                     # Loads Snap.js script
```

## Next Steps

1. **Get Midtrans Credentials**: Sign up at dashboard.midtrans.com
2. **Configure Backend**: Add keys to .env
3. **Configure Frontend**: Add client key to .env
4. **Test Checkout**: Use sandbox test cards
5. **Deploy to Production**: Switch keys and environment

## Support & Documentation

- Midtrans Documentation: https://docs.midtrans.com
- Snap.js Integration: https://docs.midtrans.com/en/snap/overview
- Webhook Documentation: https://docs.midtrans.com/en/after-payment/http-notification

---
Last Updated: 2024
SmartFarm API v1.0
