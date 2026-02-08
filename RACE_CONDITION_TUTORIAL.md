# Tutorial: Mengatasi Race Condition di E-commerce

## 📚 Daftar Isi

1. [Apa itu Race Condition?](#apa-itu-race-condition)
2. [Contoh Kasus Nyata](#contoh-kasus-nyata)
3. [Implementasi Solusi](#implementasi-solusi)
4. [Testing Race Condition](#testing-race-condition)
5. [Best Practices](#best-practices)

---

## 🎯 Apa itu Race Condition?

**Race Condition** adalah kondisi dimana hasil akhir dari suatu operasi **bergantung pada urutan atau timing** dari beberapa proses yang berjalan bersamaan (concurrent).

### 🏦 Analogi: ATM dan Saldo Bank

Bayangkan Anda punya **saldo Rp 1.000.000** di bank:

```
Waktu yang Sama (Bersamaan):
┌─────────────────────────────┐  ┌─────────────────────────────┐
│ Anda di ATM Jakarta         │  │ Pasangan di ATM Bandung     │
│ Cek saldo: Rp 1.000.000 ✓   │  │ Cek saldo: Rp 1.000.000 ✓   │
│ Tarik Rp 800.000            │  │ Tarik Rp 800.000            │
│ Saldo baru: Rp 200.000      │  │ Saldo baru: Rp 200.000      │
└─────────────────────────────┘  └─────────────────────────────┘

❌ MASALAH: Kedua transaksi berhasil!
💥 Saldo akhir: Rp 200.000 (seharusnya Rp -600.000 atau salah satu ditolak)
💸 Bank rugi Rp 600.000!
```

**Ini Race Condition!** Kedua proses "berlomba" (race) mengakses data yang sama tanpa koordinasi.

---

## 🛒 Contoh Kasus Nyata

### Skenario: Flash Sale di E-commerce

**Produk:** iPhone 15 Pro - **Stok: 1 unit** - **Harga: Rp 10.000.000**

```
Detik 12:00:00.000 - Flash Sale Dimulai!

┌─────────────────────────────┐  ┌─────────────────────────────┐
│ User A (Jakarta)            │  │ User B (Surabaya)           │
├─────────────────────────────┤  ├─────────────────────────────┤
│ 1. Klik "Beli Sekarang"     │  │ 1. Klik "Beli Sekarang"     │
│ 2. Cek stok: 1 ✓            │  │ 2. Cek stok: 1 ✓            │
│ 3. Stok cukup? YA           │  │ 3. Stok cukup? YA           │
│ 4. Kurangi stok: 1 - 1 = 0  │  │ 4. Kurangi stok: 1 - 1 = 0  │
│ 5. Buat order ✓             │  │ 5. Buat order ✓             │
│ 6. Bayar Rp 10.000.000 ✓    │  │ 6. Bayar Rp 10.000.000 ✓    │
└─────────────────────────────┘  └─────────────────────────────┘

💥 MASALAH:
- Kedua user berhasil order dan bayar
- Toko cuma punya 1 iPhone
- Salah satu user akan komplain
- Reputasi toko rusak
- Potensi refund + denda
```

### Dampak Bisnis

| Aspek | Dampak |
|-------|--------|
| **Finansial** | Kehilangan produk senilai Rp 10.000.000 |
| **Reputasi** | Review buruk, kehilangan kepercayaan |
| **Operasional** | Waktu terbuang untuk handle komplain |
| **Legal** | Potensi gugatan konsumen |

---

## 🔧 Implementasi Solusi

### Solusi 1: Database Transaction + Row-Level Locking

**Konsep:**
- Gunakan **Database Transaction** untuk memastikan operasi atomic
- Gunakan **Row-Level Locking** untuk mengunci data yang sedang diproses
- Proses lain harus **menunggu** sampai lock dilepas

#### Kode Implementasi (Go + GORM)

**File: `services/order_service.go`**

```go
package services

import (
    "errors"
    "smartfarm-api/dto"
    "smartfarm-api/models"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

type OrderService interface {
    CreateOrder(req dto.CreateOrderRequest, userID uint) (*dto.OrderResponse, error)
}

type orderService struct {
    db *gorm.DB
}

func NewOrderService(db *gorm.DB) OrderService {
    return &orderService{db: db}
}

func (s *orderService) CreateOrder(req dto.CreateOrderRequest, userID uint) (*dto.OrderResponse, error) {
    var order models.Order
    
    // 🔒 MULAI DATABASE TRANSACTION
    // Semua operasi di dalam transaction akan di-commit atau di-rollback bersamaan
    err := s.db.Transaction(func(tx *gorm.DB) error {
        
        // 🔒 STEP 1: LOCK ROW PRODUK
        // SELECT * FROM products WHERE id = ? FOR UPDATE
        // Kunci baris ini sampai transaction selesai
        var product models.Product
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&product, req.ProductID).Error; err != nil {
            return err
        }
        
        // ✅ STEP 2: CEK STOK (dengan data yang sudah di-lock)
        if product.Stock < req.Quantity {
            return errors.New("stok tidak cukup")
        }
        
        // ✅ STEP 3: KURANGI STOK (Atomic operation)
        product.Stock -= req.Quantity
        if err := tx.Save(&product).Error; err != nil {
            return err
        }
        
        // ✅ STEP 4: HITUNG TOTAL HARGA
        totalPrice := product.Price * float64(req.Quantity)
        
        // ✅ STEP 5: BUAT ORDER
        order = models.Order{
            UserID:      userID,
            TotalAmount: totalPrice,
            Status:      "pending",
        }
        if err := tx.Create(&order).Error; err != nil {
            return err
        }
        
        // ✅ STEP 6: BUAT ORDER ITEM
        orderItem := models.OrderItem{
            OrderID:   order.ID,
            ProductID: product.ID,
            Quantity:  req.Quantity,
            Price:     product.Price,
        }
        if err := tx.Create(&orderItem).Error; err != nil {
            return err
        }
        
        // ✅ Semua berhasil, commit transaction
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    // Return order response
    return &dto.OrderResponse{
        ID:          order.ID,
        TotalAmount: order.TotalAmount,
        Status:      order.Status,
    }, nil
}
```

#### Cara Kerja Step-by-Step

```
User A dan User B order bersamaan (concurrent):

┌─────────────────────────────────────────────────────────────┐
│ User A (Detik 12:00:00.000)                                 │
├─────────────────────────────────────────────────────────────┤
│ 1. Mulai transaction                                        │
│ 2. 🔒 LOCK row produk ID=1 (User B harus TUNGGU!)           │
│ 3. Cek stok: 1 ✓                                            │
│ 4. Kurangi stok: 1 - 1 = 0                                  │
│ 5. Buat order #101 ✓                                        │
│ 6. Commit transaction (SEMUA BERHASIL)                      │
│ 7. 🔓 UNLOCK row produk                                     │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ User B (Detik 12:00:00.001) - MENUNGGU UNLOCK              │
├─────────────────────────────────────────────────────────────┤
│ 1. Mulai transaction                                        │
│ 2. 🔒 LOCK row produk ID=1 (Sekarang bisa!)                 │
│ 3. Cek stok: 0 ❌                                           │
│ 4. ❌ ERROR: "stok tidak cukup"                             │
│ 5. Rollback transaction (SEMUA DIBATALKAN)                  │
│ 6. 🔓 UNLOCK row produk                                     │
│ 7. Return error ke frontend                                 │
└─────────────────────────────────────────────────────────────┘

✅ HASIL: Hanya User A yang berhasil order!
✅ User B dapat pesan error yang jelas
✅ Stok tetap akurat (0)
```

---

### Solusi 2: Optimistic Locking (Version Control)

**Konsep:**
- Tambahkan kolom `version` di tabel
- Setiap update, cek apakah version masih sama
- Kalau version berubah = ada yang update duluan

#### Database Schema

```go
type Product struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `json:"name"`
    Stock       int       `json:"stock"`
    Version     int       `json:"version"` // 👈 Tambahkan ini
    // ... fields lainnya
}
```

#### Kode Implementasi

```go
func (s *orderService) CreateOrderOptimistic(req dto.CreateOrderRequest, userID uint) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        var product models.Product
        
        // STEP 1: Ambil produk dengan version-nya
        if err := tx.First(&product, req.ProductID).Error; err != nil {
            return err
        }
        
        // STEP 2: Cek stok
        if product.Stock < req.Quantity {
            return errors.New("stok tidak cukup")
        }
        
        // STEP 3: Update dengan CEK VERSION
        // UPDATE products 
        // SET stock = stock - ?, version = version + 1 
        // WHERE id = ? AND version = ?
        result := tx.Model(&product).
            Where("id = ? AND version = ?", product.ID, product.Version).
            Updates(map[string]interface{}{
                "stock":   product.Stock - req.Quantity,
                "version": product.Version + 1,
            })
        
        // STEP 4: Cek apakah update berhasil
        if result.RowsAffected == 0 {
            // Version berubah = ada yang update duluan
            return errors.New("produk sedang diproses user lain, silakan coba lagi")
        }
        
        // STEP 5: Lanjut buat order...
        // ... (sama seperti sebelumnya)
        
        return nil
    })
}
```

---

### Perbandingan Solusi

| Aspek | Pessimistic Lock (FOR UPDATE) | Optimistic Lock (Version) |
|-------|-------------------------------|---------------------------|
| **Cara Kerja** | Lock row, tunggu selesai | Cek version saat update |
| **Performance** | Lebih lambat (ada waiting) | Lebih cepat (no waiting) |
| **Cocok untuk** | High contention (banyak conflict) | Low contention |
| **Kompleksitas** | Sederhana | Perlu retry logic |
| **Deadlock Risk** | Ada (jika lock banyak row) | Tidak ada |

**Rekomendasi untuk SmartFarm:** Gunakan **Pessimistic Lock** karena:
- Flash sale = high contention
- Lebih mudah dipahami
- Tidak perlu retry logic di frontend

---

## 🧪 Testing Race Condition

### Test 1: Simulasi Concurrent Orders (Manual)

**Tools:** Postman / cURL

```bash
# Terminal 1 - User A
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN_USER_A" \
  -d '{
    "product_id": 1,
    "quantity": 1,
    "address_id": 1
  }'

# Terminal 2 - User B (jalankan BERSAMAAN dengan Terminal 1)
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN_USER_B" \
  -d '{
    "product_id": 1,
    "quantity": 1,
    "address_id": 2
  }'
```

**Expected Result:**
- ✅ Salah satu request berhasil (200 OK)
- ✅ Salah satu request gagal (400 Bad Request: "stok tidak cukup")

---

### Test 2: Automated Test dengan Go

**File: `services/order_service_test.go`**

```go
package services

import (
    "sync"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestConcurrentOrders_RaceCondition(t *testing.T) {
    // Setup database dan service
    db := setupTestDB()
    orderService := NewOrderService(db)
    
    // Buat produk dengan stok 1
    product := models.Product{
        Name:  "iPhone 15 Pro",
        Stock: 1,
        Price: 10000000,
    }
    db.Create(&product)
    
    // Buat 2 user
    user1 := models.User{Name: "User A", Email: "a@test.com"}
    user2 := models.User{Name: "User B", Email: "b@test.com"}
    db.Create(&user1)
    db.Create(&user2)
    
    // Simulasi 2 request bersamaan
    var wg sync.WaitGroup
    var successCount int
    var mu sync.Mutex
    
    // Goroutine 1 - User A
    wg.Add(1)
    go func() {
        defer wg.Done()
        req := dto.CreateOrderRequest{
            ProductID: product.ID,
            Quantity:  1,
        }
        _, err := orderService.CreateOrder(req, user1.ID)
        if err == nil {
            mu.Lock()
            successCount++
            mu.Unlock()
        }
    }()
    
    // Goroutine 2 - User B
    wg.Add(1)
    go func() {
        defer wg.Done()
        req := dto.CreateOrderRequest{
            ProductID: product.ID,
            Quantity:  1,
        }
        _, err := orderService.CreateOrder(req, user2.ID)
        if err == nil {
            mu.Lock()
            successCount++
            mu.Unlock()
        }
    }()
    
    // Tunggu kedua goroutine selesai
    wg.Wait()
    
    // ✅ ASSERTION: Hanya 1 yang berhasil
    assert.Equal(t, 1, successCount, "Hanya 1 order yang boleh berhasil")
    
    // ✅ ASSERTION: Stok jadi 0
    var updatedProduct models.Product
    db.First(&updatedProduct, product.ID)
    assert.Equal(t, 0, updatedProduct.Stock, "Stok harus 0")
    
    // ✅ ASSERTION: Hanya ada 1 order
    var orderCount int64
    db.Model(&models.Order{}).Count(&orderCount)
    assert.Equal(t, int64(1), orderCount, "Hanya ada 1 order")
}
```

**Jalankan Test:**
```bash
cd backend-go
go test ./services -v -run TestConcurrentOrders
```

---

### Test 3: Load Testing dengan Apache Bench

```bash
# Install Apache Bench (jika belum ada)
# Windows: Download dari https://www.apachelounge.com/download/

# Jalankan 100 request concurrent
ab -n 100 -c 100 -p order.json -T application/json \
   -H "Authorization: Bearer YOUR_TOKEN" \
   http://localhost:8080/orders
```

**File `order.json`:**
```json
{
  "product_id": 1,
  "quantity": 1,
  "address_id": 1
}
```

**Expected Result:**
- ✅ Hanya 1 request berhasil (jika stok = 1)
- ✅ 99 request gagal dengan error "stok tidak cukup"

---

## 📊 Monitoring & Debugging

### 1. Tambahkan Logging

```go
func (s *orderService) CreateOrder(req dto.CreateOrderRequest, userID uint) (*dto.OrderResponse, error) {
    log.Printf("[ORDER] User %d mencoba order produk %d qty %d", userID, req.ProductID, req.Quantity)
    
    err := s.db.Transaction(func(tx *gorm.DB) error {
        var product models.Product
        
        log.Printf("[LOCK] Mencoba lock produk %d", req.ProductID)
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&product, req.ProductID).Error; err != nil {
            log.Printf("[ERROR] Gagal lock produk: %v", err)
            return err
        }
        log.Printf("[LOCK] Berhasil lock produk %d, stok: %d", req.ProductID, product.Stock)
        
        if product.Stock < req.Quantity {
            log.Printf("[REJECT] Stok tidak cukup untuk user %d", userID)
            return errors.New("stok tidak cukup")
        }
        
        product.Stock -= req.Quantity
        if err := tx.Save(&product).Error; err != nil {
            log.Printf("[ERROR] Gagal update stok: %v", err)
            return err
        }
        log.Printf("[SUCCESS] Stok dikurangi, sisa: %d", product.Stock)
        
        // ... lanjut buat order
        
        return nil
    })
    
    if err != nil {
        log.Printf("[FAILED] Order gagal untuk user %d: %v", userID, err)
        return nil, err
    }
    
    log.Printf("[SUCCESS] Order berhasil untuk user %d", userID)
    return &orderResponse, nil
}
```

### 2. Monitor Database Locks

**MySQL:**
```sql
-- Lihat transaksi yang sedang berjalan
SELECT * FROM information_schema.innodb_trx;

-- Lihat lock yang sedang aktif
SELECT * FROM information_schema.innodb_locks;

-- Lihat lock yang sedang menunggu
SELECT * FROM information_schema.innodb_lock_waits;
```

---

## ✅ Best Practices

### 1. **Selalu Gunakan Transaction untuk Operasi Kritis**

❌ **JANGAN:**
```go
// Tanpa transaction - BERBAHAYA!
product.Stock -= quantity
db.Save(&product)
db.Create(&order)
```

✅ **LAKUKAN:**
```go
// Dengan transaction - AMAN!
db.Transaction(func(tx *gorm.DB) error {
    // Semua operasi di sini atomic
    return nil
})
```

### 2. **Lock Hanya Data yang Diperlukan**

❌ **JANGAN:**
```go
// Lock seluruh tabel - LAMBAT!
tx.Exec("LOCK TABLES products WRITE")
```

✅ **LAKUKAN:**
```go
// Lock hanya 1 row - CEPAT!
tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, id)
```

### 3. **Timeout untuk Mencegah Deadlock**

```go
// Set timeout 5 detik
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // Operasi di sini
    return nil
})
```

### 4. **Retry Logic untuk Optimistic Locking**

```go
func CreateOrderWithRetry(req dto.CreateOrderRequest, userID uint, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := CreateOrderOptimistic(req, userID)
        if err == nil {
            return nil // Berhasil
        }
        if err.Error() == "version conflict" {
            time.Sleep(time.Millisecond * 100) // Tunggu sebentar
            continue // Coba lagi
        }
        return err // Error lain, langsung return
    }
    return errors.New("max retries exceeded")
}
```

### 5. **Handle Error dengan Baik di Frontend**

```typescript
// frontend/src/services/orderService.ts
export async function createOrder(data: CreateOrderRequest) {
  try {
    const response = await http.post('/orders', data)
    return { success: true, data: response.data }
  } catch (error: any) {
    if (error.response?.data?.error === 'stok tidak cukup') {
      return { 
        success: false, 
        message: 'Maaf, produk sudah habis. Silakan pilih produk lain.' 
      }
    }
    return { 
      success: false, 
      message: 'Terjadi kesalahan. Silakan coba lagi.' 
    }
  }
}
```

---

## 🎓 Kesimpulan

### Kapan Race Condition Terjadi?

Race Condition terjadi saat:
1. ✅ Ada **data yang diakses bersamaan** (concurrent access)
2. ✅ Data tersebut **diubah** (write operation)
3. ✅ **Tidak ada koordinasi** antar proses

### Solusi Race Condition

| Teknik | Kelebihan | Kekurangan |
|--------|-----------|------------|
| **Pessimistic Lock** | Simple, reliable | Slower, deadlock risk |
| **Optimistic Lock** | Fast, no deadlock | Need retry logic |
| **Queue System** | Scalable | Complex setup |

### Implementasi di SmartFarm

✅ **Sudah Diimplementasikan:**
- Database Transaction untuk atomicity
- Row-Level Locking (FOR UPDATE) untuk mencegah race condition
- Error handling yang jelas
- Logging untuk debugging

✅ **Hasil:**
- Stok produk selalu akurat
- Tidak ada overselling
- User experience yang baik (error message jelas)

---

## 📚 Referensi

1. **GORM Documentation - Transactions**  
   https://gorm.io/docs/transactions.html

2. **MySQL Locking Reads**  
   https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html

3. **Go Concurrency Patterns**  
   https://go.dev/blog/pipelines

4. **Database Isolation Levels**  
   https://en.wikipedia.org/wiki/Isolation_(database_systems)

---

**Author:** SmartFarm Development Team  
**Date:** February 2026  
**Version:** 1.0
