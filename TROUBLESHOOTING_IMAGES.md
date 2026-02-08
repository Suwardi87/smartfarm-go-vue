# Troubleshooting: Gambar Produk Tidak Muncul

## ✅ Yang Sudah Diperbaiki

1. **Backend**: Seeder sudah menggunakan 10 URL gambar Unsplash yang bervariasi
2. **Database**: Auto-cleanup produk lama sebelum seeding
3. **Frontend**: Error handling untuk gambar yang gagal load

## 🔍 Verifikasi Backend

Backend sudah mengirim URL dengan benar:

```bash
# Test API
curl http://localhost:8080/products?page=1&limit=1
```

**Expected Output:**
```json
{
  "data": [{
    "id": 4,
    "name": "Bayam Organik Segar",
    "image_url": "https://images.unsplash.com/photo-1598170845058-32b9d6a5da37?auto=format&fit=crop&q=80&w=400",
    ...
  }]
}
```

✅ URL lengkap dengan `https://images.unsplash.com/...`

## 🌐 Cara Mengatasi (Browser Cache)

### Solusi 1: Hard Refresh Browser

**Chrome/Edge:**
- Windows: `Ctrl + Shift + R` atau `Ctrl + F5`
- Mac: `Cmd + Shift + R`

**Firefox:**
- Windows: `Ctrl + Shift + R` atau `Ctrl + F5`
- Mac: `Cmd + Shift + R`

### Solusi 2: Clear Browser Cache

**Chrome/Edge:**
1. Tekan `F12` untuk buka DevTools
2. Klik kanan pada tombol Refresh
3. Pilih "Empty Cache and Hard Reload"

**Firefox:**
1. Tekan `Ctrl + Shift + Delete`
2. Pilih "Cached Web Content"
3. Klik "Clear Now"

### Solusi 3: Incognito/Private Mode

Buka aplikasi di mode incognito:
- Chrome: `Ctrl + Shift + N`
- Firefox: `Ctrl + Shift + P`
- Edge: `Ctrl + Shift + N`

## 🔧 Debugging di Browser

### 1. Cek Console Errors

1. Buka DevTools (`F12`)
2. Pilih tab "Console"
3. Refresh halaman
4. Lihat apakah ada error terkait gambar

**Expected:** Tidak ada error, atau warning "Failed to load image" jika gambar gagal

### 2. Cek Network Tab

1. Buka DevTools (`F12`)
2. Pilih tab "Network"
3. Filter: `Img` atau `All`
4. Refresh halaman
5. Lihat request ke `images.unsplash.com`

**Expected:**
- Request ke Unsplash: Status `200 OK`
- Content-Type: `image/jpeg` atau `image/webp`

### 3. Test Manual Image URL

Copy salah satu URL gambar dari API response, paste di browser:

```
https://images.unsplash.com/photo-1598170845058-32b9d6a5da37?auto=format&fit=crop&q=80&w=400
```

**Expected:** Gambar sayuran/buah muncul

## 🚨 Jika Masih Belum Muncul

### Kemungkinan 1: CORS Issue

Unsplash biasanya allow CORS, tapi jika ada masalah:

**Solusi:** Tambahkan proxy di Vite config

```typescript
// vite.config.ts
export default defineConfig({
  server: {
    proxy: {
      '/unsplash': {
        target: 'https://images.unsplash.com',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/unsplash/, '')
      }
    }
  }
})
```

### Kemungkinan 2: Firewall/Antivirus Block

Beberapa antivirus block request ke domain external.

**Solusi:**
1. Whitelist `images.unsplash.com` di antivirus
2. Atau gunakan gambar lokal (download dulu)

### Kemungkinan 3: Internet Issue

Unsplash memerlukan koneksi internet.

**Test:**
```bash
ping images.unsplash.com
```

**Expected:** Reply dari server Unsplash

## ✅ Checklist Troubleshooting

- [ ] Hard refresh browser (`Ctrl + Shift + R`)
- [ ] Clear browser cache
- [ ] Test di incognito mode
- [ ] Cek console errors (F12)
- [ ] Cek network tab untuk request Unsplash
- [ ] Test URL gambar langsung di browser
- [ ] Cek koneksi internet
- [ ] Restart dev server (frontend & backend)

## 📸 Screenshot untuk Debugging

Jika masih bermasalah, screenshot:
1. **Console tab** - untuk lihat errors
2. **Network tab** - untuk lihat request status
3. **Halaman marketplace** - untuk lihat kondisi saat ini

---

**Update Terakhir:** 6 Februari 2026  
**Status:** Gambar seharusnya sudah muncul setelah hard refresh
