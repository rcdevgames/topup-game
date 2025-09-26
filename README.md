# Topup Game Online

Aplikasi web untuk topup game online yang dibangun dengan React, Vite, TailwindCSS, dan teknologi modern lainnya.

## Teknologi yang Digunakan

- **Runtime**: Bun
- **Framework**: React 18 + Vite
- **Styling**: TailwindCSS
- **State Management**: Zustand
- **Form Handling**: React Hook Form + Yup
- **Routing**: React Router DOM
- **Notifications**: SweetAlert2

## Fitur Utama

### 🏠 Beranda
- Pencarian produk game
- Filter berdasarkan kategori
- Grid produk yang responsif (2 kolom di mobile, hingga 4 di desktop)
- Load more untuk produk tambahan

### 💳 Transaksi
- Riwayat transaksi dengan status (berhasil, pending, gagal)
- Load more untuk transaksi tambahan
- Hanya dapat diakses setelah login

### 👤 Profil
- Edit informasi pribadi (nama, nomor HP, password)
- CRUD akun game (tambah, edit, hapus)
- Logout dari aplikasi
- Hanya dapat diakses setelah login

### 🔐 Login
- Form login dengan validasi
- Demo login untuk testing
- Redirect protection untuk halaman yang memerlukan auth

## Instalasi dan Menjalankan Aplikasi

### Prerequisites
- Pastikan Bun sudah terinstall di sistem Anda

### Instalasi Dependencies
```bash
bun install
```

### Menjalankan Development Server
```bash
bun run dev
```

Aplikasi akan berjalan di http://localhost:5173/

### Build untuk Production
```bash
bun run build
```

### Preview Build
```bash
bun run preview
```

## Demo Login

Untuk testing, gunakan kredensial berikut:
- **Nomor HP**: `08123456789`
- **Password**: `password`

Atau klik tombol "Demo Login" di halaman login.

## Struktur Project

```
src/
├── components/          # Komponen reusable
│   ├── BottomNavbar.jsx
│   ├── ProductCard.jsx
│   ├── TransactionCard.jsx
│   └── GameAccountForm.jsx
├── pages/              # Halaman utama
│   ├── Home.jsx
│   ├── Transactions.jsx
│   ├── Profile.jsx
│   └── Login.jsx
├── store/              # Zustand stores
│   ├── authStore.js
│   └── gameStore.js
├── utils/              # Utilities
│   └── validation.js
├── styles/             # CSS styles
│   └── globals.css
└── App.jsx            # Main app component
```

## Fitur UI/UX

### Design System
- **Warna Utama**: Nuansa biru gelap (#35374B, #344955, #50727B)
- **Accent**: Hijau untuk success (#78A083)
- **Typography**: Inter font family
- **Shadows**: Subtle elevation untuk cards

### Responsivitas
- Mobile-first design
- Bottom navigation untuk mobile
- Grid system yang adaptif
- Touch-friendly button sizes

### Accessibility
- Focus rings untuk keyboard navigation
- Semantic HTML markup
- Aria labels untuk screen readers
- Color contrast sesuai WCAG AA

### Micro-interactions
- Hover effects pada buttons dan cards
- Smooth transitions
- Loading states
- Success/error notifications dengan SweetAlert2

## Cara Testing

### Manual Testing Flow
1. **Beranda**: 
   - Test pencarian produk
   - Filter berdasarkan kategori
   - Klik "Beli Sekarang" untuk menambah ke cart
   - Test load more

2. **Login**: 
   - Coba akses Transaksi/Profil tanpa login (akan redirect)
   - Login dengan demo credentials
   - Test form validation

3. **Transaksi**: 
   - Lihat riwayat transaksi
   - Test load more

4. **Profil**: 
   - Edit informasi pribadi
   - Tambah akun game baru
   - Edit dan hapus akun game
   - Test logout

### Responsive Testing
- Resize browser window untuk test responsivitas
- Test di mobile view (320px width)
- Pastikan bottom navigation bekerja dengan baik

## Development Notes

### State Management
- Auth state (login status, user data) disimpan di `authStore`
- Game data (products, transactions, accounts) disimpan di `gameStore`
- Local component state tetap menggunakan `useState` untuk form inputs

### Form Validation
- Menggunakan Yup schema untuk validasi
- React Hook Form untuk form handling
- Error messages dalam bahasa Indonesia

### Data Dummy
- Sample products dan transactions sudah disediakan
- Game accounts CRUD fully functional
- Search dan filter bekerja dengan data dummy

## Production Considerations

Untuk production deployment, pertimbangkan:
1. Integrasi dengan real API endpoints
2. Environment variables untuk konfigurasi
3. Error boundary components
4. Loading states yang lebih comprehensive
5. Caching strategy
6. SEO optimization
7. Performance monitoring

## Browser Support

- Chrome/Edge (terbaru)
- Firefox (terbaru)
- Safari (terbaru)
- Mobile browsers

## License

MIT License