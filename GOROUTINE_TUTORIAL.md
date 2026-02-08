# Tutorial: Goroutine - Concurrent Programming di Go

## 📚 Daftar Isi

1. [Apa itu Goroutine?](#apa-itu-goroutine)
2. [Konsep Dasar](#konsep-dasar)
3. [Implementasi di SmartFarm](#implementasi-di-smartfarm)
4. [Best Practices](#best-practices)
5. [Common Pitfalls](#common-pitfalls)

---

## 🎯 Apa itu Goroutine?

**Goroutine** adalah **lightweight thread** yang dikelola oleh Go runtime untuk menjalankan fungsi secara **concurrent** (bersamaan).

### 🏪 Analogi: Kasir Supermarket

**Sequential (Tanpa Goroutine):**
```
Kasir 1: Pelanggan A → Scan → Bayar → Selesai (2 menit)
         Pelanggan B → Scan → Bayar → Selesai (2 menit)
         Pelanggan C → Scan → Bayar → Selesai (2 menit)

Total waktu: 6 menit
```

**Concurrent (Dengan Goroutine):**
```
Kasir 1: Pelanggan A → Scan → Bayar → Selesai (2 menit)
Kasir 2: Pelanggan B → Scan → Bayar → Selesai (2 menit)  } Bersamaan!
Kasir 3: Pelanggan C → Scan → Bayar → Selesai (2 menit)  }

Total waktu: 2 menit (3x lebih cepat!)
```

---

## 📝 Konsep Dasar

### 1. Membuat Goroutine

**Syntax:** Tambahkan `go` sebelum function call

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

func main() {
    // Sequential (normal)
    sayHello("Alice")  // Tunggu selesai
    sayHello("Bob")    // Baru jalan

    // Concurrent (goroutine)
    go sayHello("Charlie")  // Langsung lanjut, tidak tunggu
    go sayHello("David")    // Jalan bersamaan dengan Charlie

    // Tunggu sebentar agar goroutine selesai
    time.Sleep(1 * time.Second)
}
```

**Output:**
```
Hello, Alice!
Hello, Bob!
Hello, Charlie!
Hello, David!
```

---

### 2. WaitGroup - Menunggu Goroutine Selesai

**Masalah:** `time.Sleep()` tidak efisien dan tidak akurat.

**Solusi:** Gunakan `sync.WaitGroup`

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()  // Panggil Done() saat function selesai
    
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)  // Tambah counter untuk setiap goroutine
        go worker(i, &wg)
    }

    wg.Wait()  // Tunggu sampai semua goroutine selesai (counter = 0)
    fmt.Println("All workers completed!")
}
```

**Output:**
```
Worker 1 starting
Worker 2 starting
Worker 3 starting
Worker 4 starting
Worker 5 starting
Worker 1 done
Worker 2 done
Worker 3 done
Worker 4 done
Worker 5 done
All workers completed!
```

**Cara Kerja WaitGroup:**
```
wg.Add(1)   → Counter: 0 → 1
wg.Add(1)   → Counter: 1 → 2
wg.Add(1)   → Counter: 2 → 3
wg.Done()   → Counter: 3 → 2
wg.Done()   → Counter: 2 → 1
wg.Done()   → Counter: 1 → 0  ✅ wg.Wait() selesai!
```

---

### 3. Mutex - Melindungi Shared Data

**Masalah:** Race Condition saat banyak goroutine akses data yang sama

```go
package main

import (
    "fmt"
    "sync"
)

var counter int

func increment(wg *sync.WaitGroup) {
    defer wg.Done()
    counter++  // ❌ RACE CONDITION!
}

func main() {
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go increment(&wg)
    }

    wg.Wait()
    fmt.Println("Counter:", counter)  // ❌ Hasilnya tidak pasti! (bisa 997, 1000, 1003, dll)
}
```

**Solusi:** Gunakan `sync.Mutex`

```go
package main

import (
    "fmt"
    "sync"
)

var (
    counter int
    mu      sync.Mutex  // Mutex untuk melindungi counter
)

func increment(wg *sync.WaitGroup) {
    defer wg.Done()
    
    mu.Lock()    // 🔒 Kunci sebelum akses
    counter++    // ✅ Aman!
    mu.Unlock()  // 🔓 Buka kunci setelah selesai
}

func main() {
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go increment(&wg)
    }

    wg.Wait()
    fmt.Println("Counter:", counter)  // ✅ Selalu 1000!
}
```

---

### 4. Channel - Komunikasi Antar Goroutine

**Channel** adalah "pipa" untuk mengirim data antar goroutine.

```go
package main

import "fmt"

func sendData(ch chan string) {
    ch <- "Hello from goroutine!"  // Kirim data ke channel
}

func main() {
    ch := make(chan string)  // Buat channel

    go sendData(ch)  // Jalankan goroutine

    msg := <-ch  // Terima data dari channel
    fmt.Println(msg)  // Output: Hello from goroutine!
}
```

**Buffered Channel:**

```go
ch := make(chan int, 3)  // Buffer size = 3

ch <- 1  // OK
ch <- 2  // OK
ch <- 3  // OK
ch <- 4  // ❌ BLOCK! (buffer penuh, tunggu ada yang baca)

fmt.Println(<-ch)  // Baca 1, sekarang buffer ada space
ch <- 4            // ✅ OK sekarang!
```

---

## 🚀 Implementasi di SmartFarm

### Contoh 1: Concurrent Order Processing

**Skenario:** Process 100 orders secara bersamaan

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Order struct {
    ID       int
    UserID   int
    Total    float64
    Status   string
}

func processOrder(order Order, wg *sync.WaitGroup) {
    defer wg.Done()
    
    fmt.Printf("[Order %d] Processing...\n", order.ID)
    
    // Simulasi proses: validasi, payment, update stock
    time.Sleep(100 * time.Millisecond)
    
    fmt.Printf("[Order %d] Completed!\n", order.ID)
}

func main() {
    orders := []Order{
        {ID: 1, UserID: 101, Total: 50000, Status: "pending"},
        {ID: 2, UserID: 102, Total: 75000, Status: "pending"},
        {ID: 3, UserID: 103, Total: 100000, Status: "pending"},
        // ... 97 orders lagi
    }

    var wg sync.WaitGroup
    start := time.Now()

    for _, order := range orders {
        wg.Add(1)
        go processOrder(order, &wg)
    }

    wg.Wait()
    
    elapsed := time.Since(start)
    fmt.Printf("\n✅ All orders processed in %v\n", elapsed)
}
```

**Hasil:**
- **Sequential:** 100 orders × 100ms = **10 detik**
- **Concurrent:** ~**100-200ms** (semua bersamaan!)

---

### Contoh 2: Bulk Seeder dengan Goroutine

**File:** `backend-go/seeders/bulk_seeder_concurrent.go`

```go
package seeders

import (
    "fmt"
    "log"
    "math/rand"
    "smartfarm-api/models"
    "sync"
    "time"

    "gorm.io/gorm"
)

func SeedBulkConcurrent(db *gorm.DB) {
    log.Println("🚀 Memulai Concurrent Seeding 100.000 data...")
    start := time.Now()

    // Cari atau buat farmer
    var farmer models.User
    if err := db.Where("role = ?", "petani").First(&farmer).Error; err != nil {
        farmer = models.User{
            Name:     "Pak Budi (Petani)",
            Email:    "budi@petani.com",
            Password: "hashed",
            Role:     "petani",
        }
        db.Create(&farmer)
    }

    categories := []string{"Vegetables", "Fruits", "Packages", "Herbs", "Hydroponics"}
    productNames := []string{"Tomat", "Bayam", "Kangkung", "Wortel", "Sawi", "Melon", "Semangka", "Cabai", "Bawang", "Selada"}
    adjectives := []string{"Segar", "Organik", "Super", "Pilihan", "Kebun", "Hidroponik", "Premium", "Manis", "Renyah", "Lokal"}

    totalProducts := 100000
    batchSize := 5000
    numBatches := totalProducts / batchSize

    var wg sync.WaitGroup
    
    // Process batches concurrently
    for i := 0; i < numBatches; i++ {
        wg.Add(1)
        
        go func(batchIndex int) {
            defer wg.Done()
            
            offset := batchIndex * batchSize
            var products []models.Product

            for j := 0; j < batchSize; j++ {
                name := fmt.Sprintf("%s %s #%d", 
                    productNames[rand.Intn(len(productNames))], 
                    adjectives[rand.Intn(len(adjectives))], 
                    offset+j)
                
                products = append(products, models.Product{
                    Name:           name,
                    Description:    fmt.Sprintf("Deskripsi produk berkualitas %s", name),
                    Price:          float64((rand.Intn(20) + 1) * 5000),
                    Stock:          rand.Intn(100) + 10,
                    Category:       categories[rand.Intn(len(categories))],
                    FarmerID:       farmer.ID,
                    ImageURL:       "https://images.unsplash.com/photo-1615485499978-f6952875f566?auto=format&fit=crop&q=80&w=400",
                    IsPreOrder:     rand.Float32() < 0.1,
                    IsSubscription: rand.Float32() < 0.05,
                    CreatedAt:      time.Now().Add(-time.Duration(rand.Intn(60)) * time.Hour * 24),
                })
            }

            // Insert batch
            if err := db.Create(&products).Error; err != nil {
                log.Printf("❌ Batch %d gagal: %v", batchIndex, err)
                return
            }

            log.Printf("✅ Batch %d/%d selesai (%d produk)", batchIndex+1, numBatches, len(products))
        }(i)
    }

    wg.Wait()

    elapsed := time.Since(start)
    log.Printf("✅ SELESAI! 100.000 data berhasil dibuat dalam waktu %v", elapsed)
}
```

**Perbandingan Performance:**

| Method | Waktu | Speedup |
|--------|-------|---------|
| Sequential | ~6 detik | 1x |
| Concurrent (20 goroutines) | ~1-2 detik | **3-6x lebih cepat!** |

---

### Contoh 3: Concurrent API Calls

**Skenario:** Fetch data dari multiple endpoints bersamaan

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "sync"
    "time"
)

type APIResponse struct {
    Endpoint string
    Data     interface{}
    Error    error
}

func fetchAPI(url string, ch chan<- APIResponse, wg *sync.WaitGroup) {
    defer wg.Done()

    resp, err := http.Get(url)
    if err != nil {
        ch <- APIResponse{Endpoint: url, Error: err}
        return
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        ch <- APIResponse{Endpoint: url, Error: err}
        return
    }

    var data interface{}
    json.Unmarshal(body, &data)

    ch <- APIResponse{Endpoint: url, Data: data}
}

func main() {
    endpoints := []string{
        "http://localhost:8080/products?page=1&limit=10",
        "http://localhost:8080/orders",
        "http://localhost:8080/users",
    }

    var wg sync.WaitGroup
    ch := make(chan APIResponse, len(endpoints))

    start := time.Now()

    // Fetch all endpoints concurrently
    for _, url := range endpoints {
        wg.Add(1)
        go fetchAPI(url, ch, &wg)
    }

    // Close channel when all goroutines done
    go func() {
        wg.Wait()
        close(ch)
    }()

    // Collect results
    for result := range ch {
        if result.Error != nil {
            fmt.Printf("❌ %s: %v\n", result.Endpoint, result.Error)
        } else {
            fmt.Printf("✅ %s: Success\n", result.Endpoint)
        }
    }

    elapsed := time.Since(start)
    fmt.Printf("\n⏱️  Total time: %v\n", elapsed)
}
```

**Hasil:**
- **Sequential:** 3 endpoints × 100ms = **300ms**
- **Concurrent:** **~100ms** (semua bersamaan!)

---

## ✅ Best Practices

### 1. **Selalu Gunakan WaitGroup**

❌ **JANGAN:**
```go
go doSomething()
time.Sleep(1 * time.Second)  // Tidak akurat!
```

✅ **LAKUKAN:**
```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    doSomething()
}()
wg.Wait()
```

---

### 2. **Protect Shared Data dengan Mutex**

❌ **JANGAN:**
```go
var counter int

go func() {
    counter++  // Race condition!
}()
```

✅ **LAKUKAN:**
```go
var (
    counter int
    mu      sync.Mutex
)

go func() {
    mu.Lock()
    counter++
    mu.Unlock()
}()
```

---

### 3. **Batasi Jumlah Goroutine**

❌ **JANGAN:**
```go
// Buat 1 juta goroutine sekaligus!
for i := 0; i < 1000000; i++ {
    go processItem(i)
}
```

✅ **LAKUKAN:** Gunakan Worker Pool

```go
func workerPool(jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
        results <- processItem(job)
    }
}

func main() {
    numWorkers := 10  // Batasi hanya 10 goroutine
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    var wg sync.WaitGroup

    // Start workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go workerPool(jobs, results, &wg)
    }

    // Send jobs
    for i := 0; i < 1000000; i++ {
        jobs <- i
    }
    close(jobs)

    wg.Wait()
    close(results)
}
```

---

### 4. **Handle Panic di Goroutine**

```go
func safeGoroutine(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Goroutine %d panicked: %v\n", id, r)
        }
    }()

    // Kode yang mungkin panic
    riskyOperation()
}
```

---

### 5. **Gunakan Context untuk Cancellation**

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d stopped\n", id)
            return
        default:
            fmt.Printf("Worker %d working...\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    for i := 1; i <= 3; i++ {
        go worker(ctx, i)
    }

    time.Sleep(3 * time.Second)
    fmt.Println("Main done")
}
```

---

## ⚠️ Common Pitfalls

### 1. **Goroutine Leak**

❌ **MASALAH:**
```go
func leak() {
    ch := make(chan int)
    go func() {
        ch <- 42  // ❌ BLOCK SELAMANYA! (tidak ada yang baca)
    }()
    // Goroutine tidak pernah selesai = MEMORY LEAK!
}
```

✅ **SOLUSI:**
```go
func noLeak() {
    ch := make(chan int, 1)  // Buffered channel
    go func() {
        ch <- 42  // ✅ OK, tidak block
    }()
}
```

---

### 2. **Closure Variable Problem**

❌ **MASALAH:**
```go
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)  // ❌ Semua print 5!
    }()
}
```

✅ **SOLUSI:**
```go
for i := 0; i < 5; i++ {
    go func(id int) {
        fmt.Println(id)  // ✅ Print 0, 1, 2, 3, 4
    }(i)
}
```

---

### 3. **Forget to Wait**

❌ **MASALAH:**
```go
func main() {
    go doSomething()
    // ❌ Program langsung exit, goroutine tidak selesai!
}
```

✅ **SOLUSI:**
```go
func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        doSomething()
    }()
    wg.Wait()  // ✅ Tunggu goroutine selesai
}
```

---

## 🎓 Kesimpulan

### Kapan Menggunakan Goroutine?

| Use Case | Cocok? | Alasan |
|----------|--------|--------|
| API Calls | ✅ | I/O bound, bisa parallel |
| Bulk Insert | ✅ | Database bisa handle concurrent |
| File Processing | ✅ | Independent tasks |
| Simple Sequential Logic | ❌ | Overhead tidak worth it |
| Database Transaction | ❌ | Harus sequential untuk consistency |

### Key Takeaways

1. **Goroutine = Lightweight Thread** - Bisa buat ribuan tanpa masalah
2. **WaitGroup = Koordinator** - Tunggu semua goroutine selesai
3. **Mutex = Kunci** - Protect shared data dari race condition
4. **Channel = Komunikasi** - Kirim data antar goroutine
5. **Worker Pool = Batasi Goroutine** - Jangan buat terlalu banyak

---

## 📚 Referensi

1. **Go by Example - Goroutines**  
   https://gobyexample.com/goroutines

2. **Effective Go - Concurrency**  
   https://go.dev/doc/effective_go#concurrency

3. **Go Concurrency Patterns**  
   https://go.dev/blog/pipelines

---

**Author:** SmartFarm Development Team  
**Date:** 6 Februari 2026  
**Version:** 1.0
