# SmartFarm Quick Start - Payment Gateway Testing

## Get Started in 5 Minutes

### Step 1: Get Free Midtrans Sandbox Credentials (2 min)
1. Go to https://dashboard.midtrans.com
2. Sign up with email (no credit card needed)
3. Verify email
4. Go to Settings → Access Keys
5. Copy **Server Key** and **Client Key**

### Step 2: Configure Backend (1 min)
Edit `smartfarm-api/.env`:
```env
MIDTRANS_SERVER_KEY=copy_server_key_here
MIDTRANS_CLIENT_KEY=copy_client_key_here
```

### Step 3: Configure Frontend (1 min)
Edit `vue-tailwind-admin-dashboard-main/.env`:
```env
VITE_MIDTRANS_CLIENT_KEY=copy_client_key_here
```

### Step 4: Start Services (1 min)

**Backend** (in `smartfarm-api/`):
```bash
go run cmd/main.go
```

**Frontend** (in `vue-tailwind-admin-dashboard-main/`):
```bash
npm run dev
```

## Test the Checkout Flow

1. **Open browser**: http://localhost:5173
2. **Sign in** (use existing account or create new)
3. **Add products** to cart from homepage
4. **Go to cart** and click "Proceed to Checkout"
5. **Select address** (or create new in Addresses page first)
6. **Click "Pay Now"**
7. **Use test card**: `4111 1111 1111 1111`
   - Exp: 12/25
   - CVV: 123
   - OTP (if prompted): any 6 digits
8. **Complete payment** in Snap modal
9. **Verify success**: Cart cleared, redirected to Orders page

## Sandbox Test Cards

| Card Type | Number | Exp | CVV |
|-----------|--------|-----|-----|
| Visa | 4111 1111 1111 1111 | 12/25 | 123 |
| MasterCard | 5555 5555 5555 4444 | 12/25 | 123 |
| Visa (3D Secure) | 4000 0000 0000 0002 | 12/25 | 123 |

Use any future expiry date and any 3-digit CVV.

## API Endpoints Quick Reference

### Create Payment
```bash
POST /payments
Authorization: Bearer {token}
Content-Type: application/json

{
  "order_id": 1,
  "address_id": 1,
  "amount": 150000
}
```

### Check Payment Status
```bash
GET /payments/orders/1
Authorization: Bearer {token}
```

## Troubleshooting Quick Fixes

| Issue | Solution |
|-------|----------|
| "Failed to load Snap script" | Check Client Key is set in .env |
| "order not found" | Make sure you created an order (items in cart) |
| "address not found" | Create an address in /addresses before checkout |
| Payment modal doesn't appear | Check browser console for errors, verify keys |
| Page shows blank | Clear browser cache, check console for errors |

## File Locations

```
Backend:
├── smartfarm-api/
│   ├── .env ← Add MIDTRANS_SERVER_KEY, MIDTRANS_CLIENT_KEY
│   ├── cmd/main.go
│   └── services/payment_service.go

Frontend:
├── vue-tailwind-admin-dashboard-main/
│   ├── .env ← Add VITE_MIDTRANS_CLIENT_KEY
│   └── src/views/Marketplace/Checkout.vue
```

## Next Steps After Testing

1. ✅ Test with sandbox credentials
2. ✅ Verify checkout flow works
3. ✅ Check orders page shows payment status
4. 📝 Get production Midtrans credentials
5. 🚀 Update keys and deploy
6. ✉️ Add email notifications (optional)

## Support

- Midtrans Docs: https://docs.midtrans.com
- This project: PAYMENT_SETUP.md and IMPLEMENTATION_SUMMARY.md

---

**Time to complete**: ~5 minutes ⚡

Good luck! 🎉
