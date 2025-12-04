# TurSchedule - Bot Telegram Jadwal Giliran

Bot Telegram yang dirancang untuk mengelola dan menampilkan jadwal giliran secara otomatis. Bot ini menggunakan Go dan menyediakan antarmuka interaktif melalui Telegram.

## 🎯 Fitur Utama

- ✅ Manajemen jadwal giliran yang fleksibel
- ✅ Penjadwalan otomatis dengan cron jobs
- ✅ Penyimpanan data lokal dengan JSON
- ✅ Antarmuka interaktif melalui Telegram
- ✅ Dukungan konfigurasi melalui environment variables

## 📋 Persyaratan Sistem

- **Go**: Versi 1.25 atau lebih baru
- **Git**: Untuk clone repository

## 🚀 Instalasi

### 1. Clone Repository

```bash
git clone <repository-url>
cd TurSchedule
```

### 2. Setup Environment Variables

Buat file `.env` di root direktori:

```env
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
DB_PATH=./data/schedules.json
LOG_LEVEL=INFO
```

**Cara mendapatkan Telegram Bot Token:**
1. Buka Telegram dan cari @BotFather
2. Gunakan command `/newbot`
3. Ikuti instruksi untuk membuat bot baru
4. Copy token yang diberikan ke dalam `TELEGRAM_BOT_TOKEN`

### 3. Install Dependencies

```bash
go mod download
```

### 4. Build Project

```bash
make build
```

Atau langsung jalankan:

```bash
go run main.go
```

## 📁 Struktur Direktori

```
TurSchedule/
├── main.go                 # Entry point aplikasi
├── go.mod                  # Go module definition
├── Makefile               # Build automation
├── README.md              # Dokumentasi ini
├── config/
│   └── config.go          # Konfigurasi aplikasi
├── internal/
│   ├── bot/
│   │   └── bot.go         # Logika bot Telegram
│   └── storage/
│       └── schedule.go    # Manajemen penyimpanan data
├── migrations/            # Database migrations (jika ada)
└── data/
    └── schedules.json     # Data jadwal (auto-generated)
```

## ⚙️ Konfigurasi

### Environment Variables

| Variable | Deskripsi | Default |
|----------|-----------|---------|
| `TELEGRAM_BOT_TOKEN` | Token bot Telegram | (required) |
| `DB_PATH` | Path file database JSON | `./data/schedules.json` |
| `LOG_LEVEL` | Level logging (INFO, DEBUG, ERROR) | `INFO` |

## 🎮 Penggunaan

### Menjalankan Bot

```bash
go run main.go
```

Bot akan mulai mendengarkan pesan Telegram dan siap menerima perintah.

### Perintah Bot (Contoh)

Bot ini mendukung berbagai perintah untuk mengelola jadwal. Perintah-perintah akan ditampilkan ketika Anda menulis `/start` di chat bot.

## 🔧 Development

### Membuat Build

```bash
make build
```

### Menjalankan dengan Watch Mode

```bash
make watch
```

### Running Tests

```bash
go test ./...
```

## 📦 Dependencies

Proyek ini menggunakan dependencies berikut:

- **telegram-bot-api** - Library resmi Telegram Bot API untuk Go
- **godotenv** - Untuk membaca file .env
- **cron** - Scheduler untuk penjadwalan otomatis

Lihat `go.mod` untuk informasi versi lengkap.

## 📝 File Data

### `data/schedules.json`

File ini menyimpan semua data jadwal dalam format JSON. File akan otomatis dibuat pada saat pertama kali bot dijalankan.

Contoh struktur:
```json
{
  "users": {
    "user_id": {
      "schedules": [
        {
          "id": 1,
          "name": "Jadwal Kerja",
          "items": ["Senin", "Selasa", "Rabu"]
        }
      ]
    }
  }
}
```

## 🐛 Troubleshooting

### Bot tidak merespons
- Pastikan `TELEGRAM_BOT_TOKEN` benar di file `.env`
- Periksa koneksi internet
- Lihat log untuk detail error

### Error: "DB_PATH tidak valid"
- Pastikan direktori `data/` ada
- Pastikan file memiliki permission yang tepat

### Error: "Token tidak ditemukan"
- Pastikan file `.env` ada di root direktori
- Periksa format dan value dari `TELEGRAM_BOT_TOKEN`

## 📧 Support

Jika menemukan bug atau ada saran, silakan buat issue di repository.

## 📄 License

[Tambahkan lisensi sesuai kebutuhan Anda]

## 🤝 Kontribusi

Kontribusi welcome! Silakan fork repository dan buat pull request untuk fitur atau bug fixes.

---

**Dibuat dengan ❤️ menggunakan Go dan Telegram Bot API**
