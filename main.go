package main

import (
	"log"

	"turschedule/config"
	"turschedule/internal/bot"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Gagal memuat config: %v\n", err)
	}

	log.Println("🚀 Memulai Schedule Bot...")

	// Create bot
	b, err := bot.NewBot(cfg.TelegramBotToken, cfg.DBPath)
	if err != nil {
		log.Fatalf("Gagal membuat bot: %v\n", err)
	}

	log.Println("✅ Bot berhasil dibuat")
	log.Println("⏳ Bot sedang mendengarkan pesan...")

	// Start listening
	if err := b.Start(); err != nil {
		log.Fatalf("Error saat menjalankan bot: %v\n", err)
	}
}
