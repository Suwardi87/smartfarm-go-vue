# Technical Documentation & Learning Guide: SmartFarm Pro

Selamat datang di dokumentasi teknis resmi SmartFarm. Dokumen ini dirancang sebagai panduan komprehensif untuk memahami sistem yang telah kita bangun, mencakup arsitektur, fitur utama, dan logika implementasi secara profesional namun mudah dipelajari.

---

## 1. Arsitektur Sistem & Tech Stack

Proyek ini dibangun menggunakan pola **Modern Monorepo-style** dengan pemisahan tegas antara layanan data dan antarmuka pengguna.

### 1.1 Program Backend (`backend-go`)
*   **Bahasa**: Go (Golang) - Dipilih karena performa tinggi dalam menangani request konkuren.
*   **Framework**: Gin Gonic - Web framework minimalis yang sangat cepat.
*   **Database**: MySQL dengan GORM (ORM) - Memudahkan pengelolaan relasi tabel produk, order, dan user tanpa harus menulis SQL mentak.
*   **Security**: JSON Web Token (JWT) - Standar industri untuk otentikasi tanpa state (stateless).

### 1.2 Program Frontend (`frontend-vue`)
*   **Bahasa**: TypeScript (Vue 3) - Memberikan keamanan tipe data agar meminimalisir error "undefined".
*   **Styling**: Tailwind CSS - Framework CSS modern untuk desain yang responsif dan premium.
*   **State Management**: Pinia - Digunakan untuk menyimpan data user yang sedang login di seluruh halaman.

---

## 2. Fitur Utama & Implementasi Teknis

### 2.1 Sistem Marketplace & Manajemen Gambar
**Logika**: Kita mengimplementasikan sistem penyimpanan gambar statis di server backend.
- **Backend**: Saat gambar diunggah, server menyimpannya di folder `uploads/products` dan mencatat path relatif (contoh: `products/foto.jpg`) di database.
- **Frontend**: Menggunakan utility `getImageUrl` untuk menggabungkan Base URL API dengan path gambar tersebut.
- **Pelajaran**: Memisahkan Base URL dari path database memungkinkan kita memindahkan server dengan mudah hanya dengan mengubah satu baris di `.env`.

### 2.2 Dashboard Analytics (Farmer)
**Logika**: Memberikan wawasan (insight) kepada petani tentang produk yang paling banyak diminati.
- **Problem**: Grafik sering crash jika data kosong.
- **Solusi**: Kita mengimplementasikan transformator data di frontend (`analyticsService.ts`) yang menjamin data selalu dalam format koordinat `x` dan `y` yang valid untuk library grafik.

### 2.3 Manajemen Produk (Full CRUD)
**Logika**: Fitur khusus Role **Farmer** untuk mengelola inventaris.
- **Create**: Setup upload multipart/form-data untuk mengirim file gambar bersama teks.
- **Edit/Update**: Implementasi form reaktif yang mengambil data awal (prefetching) sebelum dilakukan perubahan.
- **Delete**: Penanganan penghapusan data secara aman dengan konfirmasi di sisi pengguna.

### 2.4 Transaksi & Stok (Integrity Check)
**Logika**: Menjamin stok produk berkurang dengan benar saat ada pembelian.
- **Database Transaction**: Saat order dibuat, kita menjalankan "Database Transaction". Artinya, jika stok kurang atau pembayaran gagal di tengah jalan, seluruh proses dibatalkan otomatis (Rollback) sehingga data stok tidak pernah "ngaco".
- **Nullable Fields**: Mengatur `PaymentID` dan `AddressID` sebagai opsional (null) di awal agar Order bisa dibuat dulu sebelum transaksi pembayaran selesai.

### 2.5 Gateway Pembayaran (Mock Simulation)
**Logika**: Simulasi pembayaran tanpa harus memiliki akun bank asli atau API Key produksi.
- **Bypass Logic**: Sistem mendeteksi `mock-token`. Jika terdeteksi, frontend akan melewati (skip) pemanggilan library Midtrans asli dan langsung memanggil endpoint `/payments/mock-success`.
- **Backend Safety**: Endpoint mock hanya diaktifkan jika server mendeteksi kunci Midtrans yang digunakan adalah "Placeholder". Ini menjamin keamanan agar fitur mock tidak bisa disalahgunakan di server asli (production).

---

## 3. Best Practices & Standarisasi Kode

1.  **DTO (Data Transfer Object)**:
    Setiap pertukaran data antara Go dan Vue menggunakan format JSON yang seragam (`snake_case`). Di Go, kita menggunakan tag ``json:"nama_field"``.
2.  **Environment Variables**:
    Semua konfigurasi sensitif (koneksi DB, URL API) disimpan di `.env`. Hindari menulis URL langsung di dalam kode (hardcoded).
3.  **Modularitas**:
    Fungsi-fungsi pembantu seperti format mata uang (Rupiah) dikumpulkan di satu tempat (`src/utils/formatter.ts`) agar mudah diperbaiki di masa depan.

---

## 4. Cara Melanjutkan Pengembangan (Roadmap)

Jika kamu ingin belajar lebih lanjut secara manual, cobalah langkah berikut:
1.  **Tambah Fitur**: Coba tambahkan field "Kategori" di produk. Mulai dari model di Backend, lalu update form di Frontend.
2.  **Validasi**: Tambahkan pengecekan agar harga tidak boleh minus di sisi Backend.
3.  **UI/UX**: Modifikasi warna tema di `tailwind.config.ts` untuk melihat keajaiban sistem desain Tailwind.

---

"Coding adalah proses belajar yang tiada henti. Proyek SmartFarm ini adalah bukti kemampuanmu membangun sistem yang kompleks dari awal." 🚜💨
