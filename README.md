# 🤖 TurSchedule - Bot Telegram Pengingat Jadwal

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)
![Telegram](https://img.shields.io/badge/Telegram-Bot-26A5E4?style=for-the-badge&logo=telegram)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

Bot Telegram pintar untuk mengelola jadwal dan pengingat otomatis dengan notifikasi multi-level. Sempurna untuk mengingatkan aktivitas rutin, jadwal kerja, atau kegiatan pribadi Anda.

[Fitur](#-fitur-unggulan) • [Instalasi](#-instalasi-cepat) • [Penggunaan](#-cara-penggunaan) • [Dokumentasi](#-dokumentasi-lengkap)

</div>

---

## ✨ Fitur Unggulan

### 🎯 Manajemen Jadwal Lengkap
- **Tambah Jadwal** - Buat jadwal baru dengan judul, waktu, hari, dan catatan
- **Edit Jadwal** - Ubah semua detail jadwal kapan saja
- **Hapus Jadwal** - Hapus jadwal yang tidak diperlukan
- **Lihat Daftar** - Tampilkan semua jadwal Anda dengan rapi

### ⏰ Sistem Pengingat Cerdas
- **Triple Reminder** - Notifikasi 60, 30, dan 5 menit sebelum waktu
- **Notifikasi Utama** - Pesan khusus tepat pada waktu yang dijadwalkan
- **Dua Mode Pengingat**:
  - 🔔 **Sekali** - Reminder hanya dikirim satu kali, lalu jadwal dihapus otomatis
  - 🔊 **Berkali-kali** - Reminder berulang setiap minggu

### 🎨 Antarmuka User-Friendly
- Keyboard interaktif untuk kemudahan penggunaan
- Validasi input otomatis
- Pesan error yang jelas dan informatif
- Dukungan Bahasa Indonesia

### 🛡️ Keamanan & Reliabilitas
- Penyimpanan data lokal (JSON)
- Validasi judul unik per user
- Pengecekan format waktu otomatis
- Cron scheduler yang andal

---

## 📋 Persyaratan Sistem

| Requirement | Versi | Keterangan |
|-------------|-------|------------|
| **Go** | 1.25+ | Runtime utama |
| **Git** | Latest | Clone repository |
| **OS** | Linux/Mac/Windows | Cross-platform |

---

## 🚀 Instalasi Cepat

### 1️⃣ Clone Repository

```bash
git clone https://github.com/yourusername/TurSchedule.git
cd TurSchedule
```

### 2️⃣ Dapatkan Telegram Bot Token

1. Buka Telegram dan cari **@BotFather**
2. Kirim command `/newbot`
3. Ikuti instruksi (nama bot & username)
4. Copy token yang diberikan

### 3️⃣ Setup Environment Variables

Buat file `.env` di root direktori:

```env
TELEGRAM_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
DB_PATH=./data/schedules.json
LOG_LEVEL=INFO
```

### 4️⃣ Install Dependencies

```bash
go mod download
```

### 5️⃣ Jalankan Bot

```bash
# Development
go run main.go

# Production (Build terlebih dahulu)
make build
./bin/turschedule
```

✅ **Bot siap digunakan!** Buka Telegram dan mulai chat dengan bot Anda.

---

## 🎮 Cara Penggunaan

### Perintah Utama

| Perintah | Fungsi | Contoh |
|----------|--------|--------|
| `/start` | Memulai bot & lihat panduan | `/start` |
| `/add` | Tambah jadwal baru | `/add` |
| `/list` | Lihat semua jadwal | `/list` |
| `/edit` | Edit jadwal yang ada | `/edit` |
| `/delete` | Hapus jadwal | `/delete` |
| `/help` | Tampilkan bantuan | `/help` |

### 📝 Contoh Penggunaan: Membuat Jadwal

1. **Mulai Proses**
   ```
   /add
   ```

2. **Masukkan Nama Jadwal**
   ```
   Rapat Tim
   ```

3. **Pilih Waktu**
   ```
   Pilih dari keyboard: 09:00
   ```

4. **Pilih Hari**
   ```
   Pilih: Senin (Monday)
   Pilih: Rabu (Wednesday)
   Tekan: 🔄 Selesai Pilih
   ```

5. **Tambah Catatan (Opsional)**
   ```
   Ruang Meeting lantai 3
   ```

6. **Pilih Tipe Reminder**
   ```
   🔊 Berkali-kali (untuk reminder mingguan)
   ```

✅ **Jadwal berhasil dibuat!** Bot akan mengirim reminder otomatis.

### ⏰ Cara Kerja Reminder

Contoh: Jadwal **Senin 09:00**

| Waktu | Notifikasi | Pesan |
|-------|------------|-------|
| 08:00 | 🔔 Reminder 1 | "⏰ Pengingat 60 menit sebelum: Rapat Tim" |
| 08:30 | 🔔 Reminder 2 | "⏰ Pengingat 30 menit sebelum: Rapat Tim" |
| 08:55 | 🔔 Reminder 3 | "⏰ Pengingat 5 menit sebelum: Rapat Tim" |
| 09:00 | 🔴 **UTAMA** | "🔔 WAKTUNYA SEKARANG! Rapat Tim" |

---

## 📁 Struktur Project

```
TurSchedule/
├── 📄 main.go                 # Entry point aplikasi
├── 📄 go.mod                  # Go module dependencies
├── 📄 go.sum                  # Dependency checksums
├── 📄 Makefile               # Build automation scripts
├── 📄 README.md              # Dokumentasi (file ini)
├── 📄 .env                   # Environment variables (buat sendiri)
├── 📁 config/
│   └── config.go             # Load & parse konfigurasi
├── 📁 internal/
│   ├── bot/
│   │   └── bot.go            # Core bot logic & handlers
│   └── storage/
│       └── schedule.go       # JSON storage management
└── 📁 data/
    └── schedules.json        # Database jadwal (auto-generated)
```

---

## ⚙️ Konfigurasi

### Environment Variables

| Variable | Tipe | Default | Deskripsi |
|----------|------|---------|-----------|
| `TELEGRAM_BOT_TOKEN` | **Required** | - | Token dari @BotFather |
| `DB_PATH` | Optional | `./data/schedules.json` | Lokasi file database |
| `LOG_LEVEL` | Optional | `INFO` | Level logging (INFO/DEBUG/ERROR) |

### Contoh `.env`

```env
# Bot Configuration
TELEGRAM_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz

# Storage
DB_PATH=./data/schedules.json

# Logging
LOG_LEVEL=INFO
```

---

## 🗄️ Format Data

### Struktur `schedules.json`

```json
{
  "123456789": [
    {
      "id": "123456789_1701234567",
      "user_id": 123456789,
      "title": "Rapat Tim",
      "time": "09:00",
      "days": ["Monday", "Wednesday"],
      "note": "Ruang Meeting lantai 3",
      "reminder_type": "recurring",
      "reminder_times": [60, 30, 5],
      "reminder_sent": {}
    }
  ]
}
```

### Field Explanation

| Field | Tipe | Deskripsi |
|-------|------|-----------|
| `id` | string | Unique identifier (userID_timestamp) |
| `user_id` | int64 | Telegram user ID |
| `title` | string | Nama jadwal (unik per user) |
| `time` | string | Format HH:MM (24-jam) |
| `days` | []string | Array hari (Monday, Tuesday, ...) |
| `note` | string | Catatan opsional |
| `reminder_type` | string | "once" atau "recurring" |
| `reminder_times` | []int | Menit sebelum waktu (default: 60,30,5) |
| `reminder_sent` | map | Tracking reminder yang sudah terkirim |

---

## 🔧 Development

### Build Commands

```bash
# Build binary
make build

# Run development mode
make dev

# Run with auto-reload (jika ada air/realize)
make watch

# Clean build artifacts
make clean
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests dengan coverage
go test -cover ./...

# Run tests verbose
go test -v ./...
```

### Project Makefile

```makefile
.PHONY: build run clean

build:
	go build -o bin/turschedule main.go

run:
	go run main.go

clean:
	rm -rf bin/
```

---

## 📦 Dependencies

| Package | Versi | Fungsi |
|---------|-------|--------|
| `telegram-bot-api` | v5 | Telegram Bot API wrapper |
| `godotenv` | latest | Load environment variables |
| `cron` | v3 | Job scheduler |

### Install Semua Dependencies

```bash
go mod tidy
```

---

## 🐛 Troubleshooting

### ❌ Bot Tidak Merespons

**Penyebab:**
- Token bot salah atau tidak valid
- Koneksi internet bermasalah
- Bot diblokir oleh Telegram

**Solusi:**
```bash
# Cek token di .env
cat .env | grep TELEGRAM_BOT_TOKEN

# Test koneksi
curl -X GET "https://api.telegram.org/bot<TOKEN>/getMe"

# Cek log aplikasi
tail -f logs/bot.log
```

### ❌ Error: "DB_PATH tidak valid"

**Solusi:**
```bash
# Buat direktori data
mkdir -p data

# Cek permission
chmod 755 data/

# Set DB_PATH di .env
echo "DB_PATH=./data/schedules.json" >> .env
```

### ❌ Reminder Tidak Terkirim

**Penyebab:**
- Format waktu salah
- Cron expression error
- Bot restart dan cron belum reload

**Solusi:**
1. Cek format waktu (harus HH:MM)
2. Restart bot untuk reload semua cron jobs
3. Cek log untuk error cron

### ❌ Jadwal Terhapus Otomatis

**Penyebab:** Tipe reminder = "Sekali"

**Solusi:** Gunakan tipe "Berkali-kali" untuk reminder berulang

---

## 🚀 Deployment

### Deploy ke VPS/Server

```bash
# 1. Clone di server
git clone https://github.com/yourusername/TurSchedule.git
cd TurSchedule

# 2. Setup environment
cp .env.example .env
nano .env  # Edit sesuai kebutuhan

# 3. Build
go build -o turschedule main.go

# 4. Jalankan dengan systemd
sudo nano /etc/systemd/system/turschedule.service
```

**File `turschedule.service`:**

```ini
[Unit]
Description=TurSchedule Telegram Bot
After=network.target

[Service]
Type=simple
User=youruser
WorkingDirectory=/path/to/TurSchedule
ExecStart=/path/to/TurSchedule/turschedule
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
# Enable & start service
sudo systemctl enable turschedule
sudo systemctl start turschedule
sudo systemctl status turschedule
```

---

## 🤝 Kontribusi

Kontribusi sangat diterima! Berikut cara berkontribusi:

1. **Fork** repository ini
2. **Buat branch** fitur baru (`git checkout -b fitur-amazing`)
3. **Commit** perubahan (`git commit -m 'Tambah fitur amazing'`)
4. **Push** ke branch (`git push origin fitur-amazing`)
5. **Buat Pull Request**

### Guidelines

- Ikuti style guide Go standard
- Tambahkan test untuk fitur baru
- Update dokumentasi jika perlu
- Gunakan commit message yang jelas

---

## 📝 Roadmap

- [ ] Dukungan timezone dinamis
- [ ] Export/import jadwal (JSON/CSV)
- [ ] Reminder custom (atur sendiri menit sebelumnya)
- [ ] Notifikasi suara/sticker
- [ ] Web dashboard untuk monitoring
- [ ] Multi-language support
- [ ] Database PostgreSQL/MySQL option
- [ ] Docker containerization

---

## 📄 License

Proyek ini dilisensikan di bawah **MIT License** - lihat file [LICENSE](LICENSE) untuk detail.

```
MIT License

Copyright (c) 2024 TurSchedule

Permission is hereby granted, free of charge, to any person obtaining a copy...
```

---

## 👨‍💻 Author

**Fatr Rahman**

- Telegram: [@FtrRahman](https://t.me/FtrRahman)
- GitHub: [@yourusername](https://github.com/yourusername)

---

## 💖 Support

Jika proyek ini membantu Anda, berikan ⭐ di GitHub!

Atau support dengan:
- 🐛 Laporkan bug
- 💡 Berikan saran fitur
- 📖 Perbaiki dokumentasi
- ☕ [Buy me a coffee](https://www.buymeacoffee.com/yourlink)

---

<div align="center">

**Dibuat dengan ❤️ menggunakan Go & Telegram Bot API**

[⬆ Kembali ke atas](#-turschedule---bot-telegram-pengingat-jadwal)

</div>