# Log Pengerjaan Detail: Membangun Fitur SmartFarm (Step-by-Step)

Dokumen ini berisi langkah-langkah teknis yang sangat mendetail tentang apa yang telah kita kerjakan, disusun secara kronologis agar kamu bisa mengikuti prosesnya satu per satu.

---

## TAHAP 1: Perbaikan Dasar & Profil Pengguna
**Tujuan**: Menghubungkan template statis dengan data asli dari database.

1.  **Langkah 1.1: Pembersihan `.env`**
    - **File**: `frontend-vue/.env`
    - **Tindakan**: Menghapus baris yang korup dan memastikan `VITE_API_URL=http://localhost:8080` tertulis benar.
2.  **Langkah 1.2: Membuat Nama Profil Dinamis**
    - **File**: `frontend-vue/src/stores/user.ts` & `frontend-vue/src/components/layout/header/UserMenu.vue`
    - **Logika**: Mengambil data dari endpoint `/me`. Ganti teks "Musharof" menjadi `userStore.user.name`.
3.  **Langkah 1.3: Mapping Sidebar**
    - **File**: `frontend-vue/src/components/layout/AppSidebar.vue`
    - **Tindakan**: Menambah link "Dashboard Farmer" dan "Produk Saya" hanya untuk user dengan `role === 'farmer'`.

---

## TAHAP 2: Sistem Gambar & Media (Integrasi Backend-Frontend)
**Tujuan**: Menampilkan gambar produk baik dari hasil upload maupun dari seeder (Unsplash).

1.  **Langkah 2.1: Normalisasi Path di Backend**
    - **File**: `backend-go/services/product_service.go`
    - **Kode**: Ubah hasil simpan gambar agar menyertakan folder: `path := "products/" + filename`.
2.  **Langkah 2.2: Utility `getImageUrl` di Frontend**
    - **File**: `frontend-vue/src/utils/image.ts`
    - **Fungsi**:
      ```typescript
      export const getImageUrl = (path: string) => {
        if (path.startsWith('http')) return path; // Jika dari Unsplash
        return `${API_BASE_URL}/uploads/${path}`; // Jika dari server kita
      }
      ```

---

## TAHAP 3: CRUD Produk Farmer (Manajemen Jualan)
**Tujuan**: Memberikan kontrol penuh bagi petani atas produk mereka.

1.  **Langkah 3.1: Backend API**
    - **File**: `backend-go/controllers/product_controller.go`
    - **Fitur**: Tambahkan handler `GetFarmerProducts`, `UpdateProduct`, dan `DeleteProduct`.
2.  **Langkah 3.2: Tampilan List (Read & Delete)**
    - **File**: `frontend-vue/src/views/Marketplace/FarmerProductList.vue`
    - **Tindakan**: Membuat tabel produk dengan tombol `Hapus` yang memanggil `productService.deleteProduct`.
3.  **Langkah 3.3: Tampilan Edit (Update)**
    - **File**: `frontend-vue/src/views/Marketplace/EditProduct.vue`
    - **Logika**: Ambil data produk berdasarkan ID di URL, tampilkan di form, lalu kirim kembali ke backend via `PUT`.

---

## TAHAP 4: Fix Dashboard Analytics (Perbaikan Grafik)
**Tujuan**: Menampilkan grafik "Trending Products" tanpa crash.

1.  **Langkah 4.1: Penambahan Data Views**
    - **File**: `backend-go/models/product.go`
    - **Tindakan**: Tambahkan kolom `Views int` di database untuk menghitung seberapa sering produk dilihat.
2.  **Langkah 4.2: Transformasi Data Grafik**
    - **File**: `frontend-vue/src/services/analyticsService.ts`
    - **Teknik**: Mengubah array produk dari API menjadi format yang dibutuhkan oleh **ApexCharts** agar tidak muncul error `x of undefined`.

---

## TAHAP 5: Keamanan Transaksi & Sistem Pembayaran
**Tujuan**: Menangani pesanan dan pembayaran secara aman dan otomatis.

1.  **Langkah 5.1: Database Locking (Fix 500 Error)**
    - **File**: `backend-go/repositories/order_repository.go`
    - **Tindakan**: Membuat fungsi `UpdatePaymentInfo`. Menggunakan query SQL `Updates` yang spesifik agar tidak membentur lock dari transaksi lain.
2.  **Langkah 5.2: Mekanisme Mock Payment**
    - **Backend (`backend-go/services/payment_service.go`)**: Jika API key tidak ada, kirim token "mock-token-xxx".
    - **Frontend (`frontend-vue/src/services/paymentService.ts`)**: Jika token diawali kata "mock", langsung konfirmasi sukses ke backend tanpa membuka popup Midtrans.
3.  **Langkah 5.3: Notifikasi Sukses**
    - **File**: `frontend-vue/src/views/Marketplace/Checkout.vue`
    - **Tindakan**: Memanggil `cart.clearCart()` setelah pembayaran berhasil agar keranjang belanja kosong kembali.

---

## TAHAP 6: Refactoring & Pindahan Folder
**Tujuan**: Mengatur proyek agar siap dikembangkan secara profesional.

1.  **Langkah 6.1: Ganti Nama Folder**
    - Ganti `smartfarm-api` -> `backend-go`.
    - Ganti `vue-tailwind-admin-dashboard-main` -> `frontend-vue`.
2.  **Langkah 6.2: Final Check Configuration**
    - Update path di terminal `cd backend-go` dan `cd frontend-vue`.
    - Pastikan `main.go` di backend tetap menjalankan `r.Run(":8080")`.

---

### Tips Cara Belajar:
- **Ikuti Urutan**: Mulailah dari TAHAP 1. Jangan loncat ke TAHAP 5 sebelum TAHAP 1-2 selesai.
- **Baca Terminal**: Selalu perhatikan log di terminal Go. Saya sudah menambahkan banyak pesan `log.Printf` untuk membantu kamu melihat pergerakan data.
- **Console Log**: Gunakan `console.log()` di frontend untuk melihat data apa yang datang dari backend.

Selesai! Dengan mengikuti log ini, kamu baru saja mempelajari bagaimana sebuah aplikasi Fullstack kompleks dibangun dari sisi dasar hingga integrasi tingkat lanjut. 🚀🚜
