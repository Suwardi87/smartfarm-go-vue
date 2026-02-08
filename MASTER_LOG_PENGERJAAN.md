# Master Log: Riwayat Lengkap Pembangunan SmartFarm Pro

Dokumen ini mencatat seluruh perjalanan pengembangan aplikasi SmartFarm dari titik nol hingga menjadi aplikasi yang fungsional. Ini adalah rangkuman teknis yang disusun secara kronologis untuk kebutuhan belajar dan audit kode.

---

## 🏗️ FASE 1: Inisialisasi & Arsitektur Dasar
*Pada awalnya, proyek dimulai dengan dua fondasi utama yang bekerja secara terpisah namun saling terhubung.*

1.  **Backend (Go - Gin Framework)**:
    - Membuat struktur folder untuk memisahkan Controller, Service, Model, dan Repository (Pattern Clean Architecture).
    - Setup Database MySQL menggunakan GORM untuk sinkronisasi tabel secara otomatis (Auto-migrate).
2.  **Frontend (Vue 3 - Tailwind CSS)**:
    - Pemilihan template dashboard modern yang responsif.
    - Setup State Management (Pinia) untuk menyimpan data user secara global.

---

## 🔑 FASE 2: Otentikasi & Identitas Pengguna (Identity System)
*Langkah untuk mengubah template statis menjadi aplikasi yang mengenal siapa penggunanya.*

1.  **Integrasi JWT (JSON Web Token)**:
    - Backend memberikan token saat login.
    - Frontend menyimpan token di `localStorage` dan mengirimnya kembali di setiap request (Auth Header).
2.  **Mapping Data Dinamis**:
    - **Masalah**: Nama user masih hardcoded ("Musharof").
    - **Solusi**: Membuat endpoint `/me` di Go. Frontend memanggilnya saat pertama kali buka web dan mengisi `userStore` dengan data asli (Nama, Role, Email).

---

## 📷 FASE 3: Revolusi Penanganan Gambar (Image Engine)
*Tahap krusial agar gambar produk muncul secara konsisten di semua halaman.*

1.  **Standarisasi Backend**:
    - Memastikan setiap produk yang disimpan memiliki path relatif: `products/nama_file.jpg`.
2.  **Frontend Image Processor**:
    - Membuat utility `getImageUrl`. Logikanya: Jika URL diawali "http" (dari seeder), gunakan langsung. Jika tidak, tambahkan prefix URL API kita (`http://localhost:8080/uploads/`).
    - **Hasil**: Gambar marketplace, keranjang, dan order history muncul 100% sukses.

---

## 🛒 FASE 4: Marketplace & Logika Logistik
*Membangun alur belanja dari pemilihan produk hingga checkout.*

1.  **Address Management**:
    - Membuat tabel Alamat User agar pembeli bisa menyimpan banyak lokasi pengiriman.
2.  **Checkout Logic**:
    - **Stok Integrity**: Menambahkan sistem "Transaction" di Database. Jika stok tinggal 1 dan ada 2 orang beli barengan, DB akan mengunci (Lock) agar stok tidak menjadi minus.
3.  **Order Persistence**:
    - Mencatat data produk yang dibeli ke dalam `OrderItems` agar meskipun harga produk asli berubah di masa depan, histori belanja user tetap akurat.

---

## 📊 FASE 5: Pengalaman Petani (Farmer Experience)
*Memberikan alat bantu bagi petani (Farmer) untuk mengelola bisnis mereka.*

1.  **Farmer Dashboard Analytics**:
    - **Fix ApexCharts**: Memperbaiki error crash grafik dengan melakukan normalisasi data API (mengubah array mentah menjadi koordinat JSON).
    - **Trending Data**: Menambahkan fitur "View Count" agar grafik menunjukkan produk mana yang paling banyak dilihat pembeli.
2.  **Full CRUD Management**:
    - Membuat halaman "Produk Saya" khusus Farmer.
    - Implementasi fitur **Tambah, Edit, dan Hapus** produk lengkap dengan upload gambar baru.

---

## 💳 FASE 6: Sistem Pembayaran (Payment Gateway)
*Bagian tersulit: Menghubungkan aplikasi dengan sistem bank/E-Wallet.*

1.  **Midtrans Integration**:
    - Menghubungkan Backend dengan Midtrans Snap API.
2.  **Fix 500 Internal Error**:
    - Menemukan BUG di mana database menolak update status order karena konflik relasi.
    - **Solusi**: Membuat fungsi repositori khusus `UpdatePaymentInfo` yang bekerja secara independen.
3.  **The "Functional Mock Mode"**:
    - Menciptakan fitur bypass di mana jika developer tidak punya API Key Midtrans, aplikasi tetap bisa mensimulasikan pembayaran sukses secara otomatis untuk testing.

---

## 📂 FASE 7: Restrukturisasi & Penyempurnaan (The Great Rename)
*Merapikan "rumah" sebelum proyek dinyatakan selesai.*

1.  **Folder Renaming**:
    - Mengubah nama folder agar lebih deskriptif:
      - `smartfarm-api` ➔ `backend-go`
      - `vue-tailwind-admin-dashboard-main` ➔ `frontend-vue`
2.  **Dokumentasi**:
    - Membuat panduan belajar mandiri (TUTORIAL_LEARNING.md dan LOG_PENGERJAAN.md).

---

## 🎨 FASE 8: Integrasi Template TailAdmin (Farmer Premium Dashboard)
*Meningkatkan estetika dan fungsionalitas dashboard petani menjadi standar profesional.*

1.  **Dashboard Refactoring**:
    - Merombak total `FarmerDashboard.vue` menggunakan layout `AdminLayout` dan komponen TailAdmin.
    - Menghubungkan metrik riil: Total Pendapatan (Revenue), Total Pesanan (Orders), Jumlah Pelanggan Unik, dan Total Produk.
2.  **Smart Insights**:
    - Menampilkan "Smart Prediction" yang mengambil data produk trending dari backend untuk membantu petani menentukan komoditas tanam yang paling diminati.
3.  **Dynamic Order History**:
    - Mengintegrasikan tabel pesanan terbaru yang secara otomatis menarik data transaksi produk milik petani tersebut.

---

---

## 🔍 FASE 9: Sistem Pencarian & Filter Katalog (Discovery Engine)
*Memudahkan pembeli menemukan produk dengan pencarian keyword dan kategori.*

1.  **Backend Dynamic Search**:
    - Memodifikasi `ProductRepository` agar query pencarian menjadi dinamis menggunakan operator `LIKE` pada SQL.
    - Mendukung pencarian berdasarkan Nama Produk, Deskripsi, maupun Kategori.
2.  **Frontend Search Integration**:
    - Menghubungkan Search Bar di Header ke sistem rute Marketplace.
    - Implementasi *Watcher* di Vue untuk memicu pembaruan data setiap kali keyword pencarian atau filter kategori berubah.
3.  **Enhanced Navigation**:
    - Mengaktifkan link "Lihat Semua" di berbagai section untuk memberikan akses penuh ke katalog produk bagi pembeli.

---

## 💡 Pelajaran Utama (Key Takeaways)
1.  **Komunikasi Data**: Selalu gunakan DTO yang sama antara Go dan Vue.
2.  **Environment**: Jangan pernah hardcode URL. Gunakan .env agar aplikasi fleksibel.
3.  **Database**: Gunakan GORM Hooks/Transaction untuk menjaga integritas data yang sensitif.

🚀 **SELESAI! Proyek SmartFarm Pro kini siap untuk tahap produksi.**
