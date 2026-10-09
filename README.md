# Gadget Marketplace & Rental API

API Backend RESTful untuk platform jual-beli gadget. Dibuat menggunakan Go (Golang) dengan framework Echo v4, terintegrasi dengan Supabase PostgreSQL, Google Gemini AI, dan Mailjet Email Relay.

---

## Fitur Utama

- **Autentikasi & User Management**
  - Registrasi & Login pengguna berbasis JWT (JSON Web Token).
  - Fitur isi ulang saldo (Top-Up Saldo).
  - Manajemen profil pengguna.

- **Katalog Gadget & Sinkronisasi**
  - Pilihan beli gadget.
  - Fitur sinkronisasi produk dari 3rd Party API (DummyJSON).
  - Pencarian, filter kategori, dan pagination produk.

- **Checkout & Pembayaran Safe-Guard**
  - Transaksi atomic berbasis Database Transaction (ACID).
  - Pengecekan kecukupan saldo dan ketersediaan stok produk otomatis.

- **Notifikasi Email Invoice**
  - Pengiriman konfirmasi transaksi dan detail invoice otomatis ke email pengguna via Mailjet SMTP Relay.

- **Asisten AI Rekomendasi Gadget**
  - Integrasi dengan Google Gemini AI (v1beta).
  - Memberikan rekomendasi gadget yang dipersonalisasi berdasarkan query/kebutuhan pengguna dan stok produk yang tersedia di database.

- **Dokumentasi API Interaktif (Swagger UI)**
  - Dokumentasi OpenAPI 2.0 terintegrasi yang dapat diuji langsung dari browser.

---

## Tech Stack & Library

- **Bahasa Pemrograman**: Go (Golang) 1.22+
- **Web Framework**: Echo v4
- **ORM & Database**: GORM & PostgreSQL (Supabase Cloud)
- **Email Service**: Mailjet SMTP Relay
- **AI Integration**: Google Gemini AI SDK (google.golang.org/genai)
- **API Documentation**: Swagger / echo-swagger
- **Testing**: Testify (stretchr/testify) & Mockery
- **Deployment**: Railway / Docker

---

## Struktur Proyek (Clean Architecture)

```
.
├── config/         # Konfigurasi aplikasi & database Supabase
├── docs/           # Berkas OpenAPI/Swagger spec & handler UI
├── handler/        # Controller / HTTP Request handlers
├── middleware/     # JWT Auth & Logger middlewares
├── models/         # Entity models & Data Transfer Objects (DTO)
├── pkg/            # Third-party utilities (Mailjet Email)
├── repository/     # Data access layer (GORM queries)
├── router/         # Routing definition & middleware assignment
├── service/        # Business logic layer & Gemini AI integration
├── ddl.sql         # Skrip DDL Database (Tabel, Foreign Key, Enum)
├── Dockerfile      # Konfigurasi containerization
├── main.go         # Entry point aplikasi
└── README.md
```

---

## Persiapan & Instalasi Lokal

### 1. Prasyarat
- Go (v1.22 atau lebih baru)
- PostgreSQL database (lokal atau akun Supabase)
- API Key Google Gemini AI & Mailjet SMTP

### 2. Clone Repository
```bash
git clone https://github.com/Jetstreck/gadget-marketplace.git
cd gadget-marketplace
```

### 3. Konfigurasi Environment Variable
Salin file `.env.example` menjadi `.env` dan isi sesuai kredensial Anda:
```bash
cp .env.example .env
```

Contoh isi `.env`:
```env
PORT=8080
ENV=development

# Database Configuration (Supabase / Postgres)
DB_DRIVER=postgres
DB_HOST=aws-0-ap-northeast-1.pooler.supabase.com
DB_PORT=6543
DB_USER=postgres.xxx
DB_PASSWORD=your_password
DB_NAME=postgres
DB_SSLMODE=require

# JWT Configuration
JWT_SECRET=super_secret_jwt_key
JWT_EXPIRATION_HOURS=24

# Third Party APIs
DUMMYJSON_URL=https://dummyjson.com

# SMTP Mailjet Configuration
SMTP_HOST=in-v3.mailjet.com
SMTP_PORT=587
SMTP_USER=your_mailjet_api_key
SMTP_PASSWORD=your_mailjet_secret_key
SENDER_EMAIL=your_verified_email@gmail.com

# AI Assistant Configuration
GEMINI_API_KEY=your_google_gemini_api_key
```

### 4. Eksekusi DDL Database
Jalankan file `ddl.sql` atau `schema.sql` pada PostgreSQL / Supabase SQL Editor Anda untuk membuat tabel-tabel yang dibutuhkan.

### 5. Jalankan Aplikasi
```bash
go run main.go
```
Aplikasi akan berjalan di `http://localhost:8080`.

---

## Menjalankan Unit Test

Untuk menjalankan seluruh pengujian unit test dan memastikan skenario logic berjalan lancar:
```bash
go test ./... -v
```

---

## Dokumentasi API

Setelah server berjalan, Anda dapat mengakses Swagger UI di:
`http://localhost:8080/swagger/index.html`

---

## Endpoint Utama

| Method | Endpoint | Deskripsi | Auth |
| :--- | :--- | :--- | :---: |
| `POST` | `/api/v1/users/register` | Registrasi user baru | Tidak |
| `POST` | `/api/v1/users/login` | Login user & dapatkan token JWT | Tidak |
| `GET` | `/api/v1/users/profile` | Lihat informasi profil user | Ya |
| `POST` | `/api/v1/users/topup` | Isi ulang saldo user | Ya |
| `GET` | `/api/v1/products` | Dapatkan daftar produk gadget | Tidak |
| `POST` | `/api/v1/products/sync` | Sync produk dari DummyJSON API | Ya |
| `POST` | `/api/v1/orders/checkout` | Transaksi beli / sewa gadget | Ya |
| `GET` | `/api/v1/orders/my-orders` | Riwayat pesanan user | Ya |
| `POST` | `/api/v1/ai/recommend` | Rekomendasi gadget berbasis AI Gemini | Ya |

---
