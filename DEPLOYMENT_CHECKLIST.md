# SmartFarm Payment Implementation - Deployment Checklist

## Pre-Deployment Verification

### Backend Setup ✓
- [ ] **Go dependencies installed**
  ```bash
  cd smartfarm-api
  go mod download
  ```

- [ ] **Database configured**
  - [ ] MySQL running on localhost:3306
  - [ ] Database `smartfarm` exists
  - [ ] GORM auto-migration will create tables on first run

- [ ] **Environment variables set** in `smartfarm-api/.env`
  - [ ] DB_HOST=localhost
  - [ ] DB_PORT=3306
  - [ ] DB_USER=root
  - [ ] DB_PASSWORD=(as configured)
  - [ ] DB_NAME=smartfarm
  - [ ] JWT_SECRET=supersecretkey
  - [ ] MIDTRANS_SERVER_KEY=your_key
  - [ ] MIDTRANS_CLIENT_KEY=your_key

- [ ] **Backend builds successfully**
  ```bash
  cd smartfarm-api
  go build -o main cmd/main.go
  ```

- [ ] **Backend runs without errors**
  ```bash
  go run cmd/main.go
  # Should show: listening on :8080
  ```

### Frontend Setup ✓
- [ ] **Node modules installed**
  ```bash
  cd vue-tailwind-admin-dashboard-main
  npm install
  ```

- [ ] **Environment variables set** in `vue-tailwind-admin-dashboard-main/.env`
  - [ ] VITE_API_URL=http://localhost:8080
  - [ ] VITE_MIDTRANS_CLIENT_KEY=your_client_key

- [ ] **Frontend builds successfully**
  ```bash
  npm run build
  # Check dist/ folder created
  ```

- [ ] **Dev server runs**
  ```bash
  npm run dev
  # Should show: Local: http://localhost:5173
  ```

### Midtrans Account Setup ✓
- [ ] **Account created** at https://dashboard.midtrans.com
- [ ] **Account verified** (email confirmation)
- [ ] **Server Key obtained** from Settings → Access Keys
- [ ] **Client Key obtained** from Settings → Access Keys
- [ ] **Environment confirmed** as Sandbox (for testing)
- [ ] **Notification URL configured** (for production: https://yourdomain.com/payments/webhook)

## Feature Verification

### User Authentication ✓
- [ ] **Sign up works**
  - Create new account with email/password
  - Verify account created in database

- [ ] **Sign in works**
  - Log in with created account
  - Verify JWT token received and stored

- [ ] **Protected routes work**
  - Accessing /checkout without auth redirects to signin
  - After signin, can access all protected routes

### Address Management ✓
- [ ] **Create address**
  - Navigate to /addresses
  - Create new address with valid fields
  - Address appears in list

- [ ] **Edit address**
  - Click edit on existing address
  - Modify fields and save
  - Changes reflected in list

- [ ] **Delete address**
  - Delete an address
  - Confirm it's removed from list

- [ ] **Set default address**
  - Mark address as default
  - See default badge appears
  - Verify only one address is default

### Cart & Products ✓
- [ ] **Add to cart**
  - Navigate to /
  - Click "Add to Cart" on a product
  - See toast notification
  - Item appears in cart

- [ ] **Cart display**
  - Navigate to /cart
  - See all items with correct details
  - Quantities match what was added

- [ ] **Update quantity**
  - Change quantity in cart
  - See total update correctly

- [ ] **Remove from cart**
  - Remove item from cart
  - Item disappears from display
  - Total updates

### Checkout Flow ✓
- [ ] **Navigate to checkout**
  - Click "Proceed to Checkout" from cart
  - Route to /checkout successful
  - Page loads without errors

- [ ] **Order summary displays**
  - See all items with correct prices
  - Subtotal calculation correct
  - Tax calculation correct (10%)
  - Total correct (subtotal + tax)

- [ ] **Address selection**
  - Default address auto-selected
  - Can select different address
  - Selected address highlighted

- [ ] **Payment initiation**
  - Click "Pay Now"
  - See loading state
  - Midtrans Snap modal appears

### Payment Processing ✓
- [ ] **Midtrans Snap loads**
  - Payment modal displays
  - No JavaScript errors in console

- [ ] **Test payment succeeds**
  - Use test card: 4111 1111 1111 1111
  - Select payment method
  - Complete payment
  - See success message

- [ ] **Payment verification**
  - Check /orders page
  - New order visible with "paid" status
  - Order shows correct amount

- [ ] **Cart cleared after payment**
  - Cart is empty after successful payment
  - Redirected to /orders automatically

### Order History ✓
- [ ] **Orders page loads**
  - Navigate to /orders
  - See list of user's orders

- [ ] **Order details visible**
  - See order items, quantities, prices
  - See order status
  - See creation date

- [ ] **Filter by status**
  - Switch between Orders and Subscriptions tabs
  - See appropriate items in each tab

## API Endpoint Testing

### Using cURL or Postman

#### 1. Register
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'
```
Expected: 201 Created with user data

#### 2. Sign In
```bash
curl -X POST http://localhost:8080/signin \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```
Expected: 200 OK with token

#### 3. Get Current User
```bash
curl -X GET http://localhost:8080/me \
  -H "Authorization: Bearer {TOKEN}"
```
Expected: 200 OK with user data

#### 4. Create Address
```bash
curl -X POST http://localhost:8080/addresses \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "label":"Home",
    "recipient_name":"John Doe",
    "phone_number":"081234567890",
    "street":"Jl. Example 123",
    "city":"Jakarta",
    "province":"DKI Jakarta",
    "postal_code":"12345"
  }'
```
Expected: 201 Created with address data

#### 5. Get My Addresses
```bash
curl -X GET http://localhost:8080/addresses \
  -H "Authorization: Bearer {TOKEN}"
```
Expected: 200 OK with array of addresses

#### 6. Create Order
```bash
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "items":[
      {"product_id":1,"quantity":2}
    ],
    "address_id":1
  }'
```
Expected: 201 Created with order data including ID

#### 7. Create Payment
```bash
curl -X POST http://localhost:8080/payments \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id":1,
    "address_id":1,
    "amount":150000
  }'
```
Expected: 201 Created with snap_token and payment_id

#### 8. Get Payment Status
```bash
curl -X GET http://localhost:8080/payments/orders/1 \
  -H "Authorization: Bearer {TOKEN}"
```
Expected: 200 OK with payment status

## Performance Checklist

- [ ] **Frontend**
  - [ ] Page loads in < 2 seconds
  - [ ] Checkout modal appears in < 1 second
  - [ ] No console errors
  - [ ] No memory leaks (check DevTools)

- [ ] **Backend**
  - [ ] API responses in < 500ms
  - [ ] Database queries optimized
  - [ ] No unhandled errors in logs
  - [ ] Payment service initializes correctly

## Security Checklist

- [ ] **Authentication**
  - [ ] JWT tokens properly validated
  - [ ] Expired tokens rejected
  - [ ] Protected routes require auth

- [ ] **Authorization**
  - [ ] Users can only access their own orders
  - [ ] Users can only use their own addresses
  - [ ] Users can only see their own payments

- [ ] **Data Protection**
  - [ ] Passwords hashed in database
  - [ ] Sensitive data not logged
  - [ ] CORS properly configured (only localhost:5173)
  - [ ] Server Key never exposed to frontend

- [ ] **Payment Security**
  - [ ] Order belongs to user before creating payment
  - [ ] Address belongs to user before payment
  - [ ] Amount matches order total
  - [ ] Midtrans handles card data (PCI compliant)

## Browser Compatibility

- [ ] **Chrome** (latest)
  - [ ] All features work
  - [ ] No console errors
  - [ ] Payment modal displays correctly

- [ ] **Firefox** (latest)
  - [ ] All features work
  - [ ] No console errors

- [ ] **Safari** (latest)
  - [ ] All features work
  - [ ] No console errors

- [ ] **Edge** (latest)
  - [ ] All features work
  - [ ] No console errors

## Mobile Responsiveness

- [ ] **Checkout page responsive**
  - [ ] Mobile view < 600px width
  - [ ] Tablet view 600px - 900px
  - [ ] Desktop view > 900px
  - [ ] All elements visible and clickable

- [ ] **Touch interactions work**
  - [ ] Buttons clickable on mobile
  - [ ] Forms editable on mobile
  - [ ] Modal scrollable if needed

## Documentation Verification

- [ ] **QUICK_START.md**
  - [ ] Instructions clear and accurate
  - [ ] All 5 steps tested and work
  - [ ] Test cards listed and working

- [ ] **PAYMENT_SETUP.md**
  - [ ] Setup instructions complete
  - [ ] API endpoints documented
  - [ ] Troubleshooting covers common issues
  - [ ] Webhook configuration explained

- [ ] **IMPLEMENTATION_SUMMARY.md**
  - [ ] All files listed
  - [ ] Features accurately described
  - [ ] Flow diagram makes sense
  - [ ] Testing checklist valid

## Deployment Readiness

### For Staging/Testing
- [ ] All checkboxes above are checked
- [ ] No critical errors in logs
- [ ] Sandbox credentials working
- [ ] All test cases passed

### For Production
- [ ] Production Midtrans credentials obtained
- [ ] HTTPS enabled on backend
- [ ] Backend deployed to production server
- [ ] Frontend deployed to CDN or server
- [ ] Environment variables updated with production keys
- [ ] Database backups configured
- [ ] Monitoring/logging configured
- [ ] Error tracking configured (e.g., Sentry)
- [ ] Email notifications configured
- [ ] Customer support process defined

## Known Limitations (to Address Later)

- [ ] Search and filter not yet implemented
- [ ] Email notifications not sent (future enhancement)
- [ ] Shipping integration pending
- [ ] Admin dashboard not available
- [ ] Refund processing manual (future automation)
- [ ] No payment retry logic (manual retry)

## Support & Next Steps

### If Issues Found
1. Check relevant error logs
2. Consult PAYMENT_SETUP.md troubleshooting
3. Check browser console for frontend errors
4. Check backend logs for API errors
5. Verify Midtrans credentials

### When Ready to Deploy
1. Get production Midtrans credentials
2. Update .env files with production keys
3. Switch payment_service.go to Production environment
4. Update webhook URL in Midtrans dashboard
5. Deploy backend and frontend
6. Run smoke tests on production
7. Monitor for errors in first 24 hours

### Ongoing Maintenance
- Monitor Midtrans dashboard for payment trends
- Review transaction logs regularly
- Handle failed payments and refunds
- Update dependencies quarterly
- Security patches immediately
- User support for payment issues

---

**Completion Date**: _______________

**Tester Name**: _______________

**Sign-off**: _______________

All checklist items verified and implementation ready for deployment! ✅
