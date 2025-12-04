package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"turschedule/internal/storage"
)

type Bot struct {
	api       *tgbotapi.BotAPI
	storage   *storage.UserSchedules
	cron      *cron.Cron
	userState map[int64]UserState
}

type UserState struct {
	Action string
	Data   map[string]interface{}
}

func NewBot(token string, dbPath string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat bot API: %w", err)
	}

	stor, err := storage.NewUserSchedules(dbPath)
	if err != nil {
		return nil, fmt.Errorf("gagal menginisialisasi storage: %w", err)
	}

	bot := &Bot{
		api:       api,
		storage:   stor,
		cron:      cron.New(),
		userState: make(map[int64]UserState),
	}

	log.Printf("Bot %s sudah aktif\n", api.Self.UserName)
	return bot, nil
}

func (b *Bot) Start() error {
	// Start cron scheduler
	b.cron.Start()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.Chat.ID
		text := update.Message.Text

		if strings.HasPrefix(text, "/") {
			b.handleCommand(userID, text)
		} else {
			b.handleMessage(userID, text)
		}
	}

	return nil
}

func (b *Bot) handleCommand(userID int64, command string) {
	parts := strings.Fields(command)
	cmd := parts[0]

	switch cmd {
	case "/start":
		b.sendMessage(userID, getHelpText())

	case "/add":
		b.sendMessageWithKeyboard(userID, "Masukkan nama jadwal:", getSkipKeyboard())
		b.userState[userID] = UserState{
			Action: "add_title",
			Data:   make(map[string]interface{}),
		}

	case "/list":
		b.listSchedules(userID)

	case "/edit":
		schedules := b.storage.GetUserSchedules(userID)
		if len(schedules) == 0 {
			b.sendMessage(userID, "Anda belum memiliki jadwal. Gunakan /add untuk membuat jadwal baru.")
			return
		}
		
		b.listSchedules(userID)
		b.sendMessage(userID, "Masukkan judul jadwal yang ingin diubah:")
		b.userState[userID] = UserState{
			Action: "edit_title",
			Data:   make(map[string]interface{}),
		}

	case "/delete":
		schedules := b.storage.GetUserSchedules(userID)
		if len(schedules) == 0 {
			b.sendMessage(userID, "Anda belum memiliki jadwal. Gunakan /add untuk membuat jadwal baru.")
			return
		}
		
		b.listSchedules(userID)
		b.sendMessage(userID, "Masukkan judul jadwal yang ingin dihapus:")
		b.userState[userID] = UserState{
			Action: "delete_title",
			Data:   make(map[string]interface{}),
		}

	case "/help":
		b.sendMessage(userID, getHelpText())

	default:
		b.sendMessage(userID, "Perintah tidak dikenal. Ketik /help untuk bantuan.")
	}
}

func (b *Bot) handleMessage(userID int64, text string) {
	state, exists := b.userState[userID]
	if !exists {
		b.sendMessage(userID, "Ketik /help untuk bantuan.")
		return
	}

	// Handle cancel button
	if text == "❌ Batal" {
		delete(b.userState, userID)
		b.sendMessage(userID, "Dibatalkan. Ketik /help untuk bantuan.")
		return
	}

	switch state.Action {
	case "add_title":
		state.Data["title"] = text
		
		// Check if title already exists for this user
		if b.storage.IsTitleExists(userID, text) {
			b.sendReplyMessage(userID, "❌ Judul sudah ada. Gunakan judul yang berbeda.")
			b.sendMessageWithKeyboard(userID, "Masukkan nama jadwal:", getSkipKeyboard())
			return
		}
		
		state.Action = "add_time"
		b.userState[userID] = state
		b.sendMessageWithKeyboard(userID, "Pilih waktu:", getTimeKeyboard())

	case "add_time":
		if !isValidTime(text) {
			b.sendReplyMessage(userID, "Format waktu tidak valid. Gunakan HH:MM (contoh: 09:30)")
			return
		}
		state.Data["time"] = text
		state.Action = "add_days"
		b.userState[userID] = state
		b.sendMessageWithKeyboard(userID, "Pilih hari (bisa pilih lebih dari satu):", getDaysKeyboard())

	case "add_days":
		// Initialize days list if not exists
		var selectedDays []string
		if d, exists := state.Data["selectedDays"]; exists {
			selectedDays = d.([]string)
		}

		// Handle "Selesai Pilih" button
		if text == "🔄 Selesai Pilih" || text == "selesai pilih" || text == "Selesai Pilih" {
			if len(selectedDays) == 0 {
				b.sendReplyMessage(userID, "Pilih minimal satu hari!")
				b.sendMessageWithKeyboard(userID, "Pilih hari (bisa pilih lebih dari satu):", getDaysKeyboard())
				return
			}
			state.Data["days"] = selectedDays
			state.Data["selectedDays"] = nil // Reset selected days
			state.Action = "add_note"
			b.userState[userID] = state
			b.sendMessageWithKeyboard(userID, "Masukkan catatan (opsional, atau ketik '-'):", getNoteKeyboard())
			return
		}

		// Parse selected day
		days := parsedays(text)
		if len(days) == 0 {
			b.sendReplyMessage(userID, "Format hari tidak valid. Pilih dari tombol yang tersedia.")
			return
		}

		// Add to selected days
		newDay := days[0]
		for _, d := range selectedDays {
			if d == newDay {
				b.sendReplyMessage(userID, "Hari sudah dipilih. Pilih hari lain atau tekan 🔄 Selesai Pilih")
				b.sendMessageWithKeyboard(userID, "Hari yang dipilih: "+strings.Join(selectedDays, ", "), getDaysKeyboard())
				return
			}
		}

		selectedDays = append(selectedDays, newDay)
		state.Data["selectedDays"] = selectedDays
		b.userState[userID] = state
		b.sendMessageWithKeyboard(userID, "✅ "+newDay+" dipilih\n\nHari yang dipilih: "+strings.Join(selectedDays, ", "), getDaysKeyboard())

	case "add_note":
		note := text
		if note == "-" || note == "Tidak ada catatan" {
			note = ""
		}
		state.Data["note"] = note
		state.Action = "add_reminder_type"
		b.userState[userID] = state
		b.sendMessageWithKeyboard(userID, "Pilih tipe reminder:", getReminderTypeKeyboard())

	case "add_reminder_type":
		reminderType := strings.ToLower(text)
		if !contains([]string{"🔔 Sekali", "🔊 Berkali-kali", "sekali", "berkali-kali"}, text) {
			b.sendReplyMessage(userID, "Pilihan tidak valid. Pilih dari tombol yang tersedia.")
			return
		}

		if reminderType == "🔔 sekali" || reminderType == "sekali" {
			state.Data["reminderType"] = "once"
		} else if reminderType == "🔊 berkali-kali" || reminderType == "berkali-kali" {
			state.Data["reminderType"] = "recurring"
		}

		// Create schedule with reminder settings
		schedule := &storage.Schedule{
			ID:            fmt.Sprintf("%d_%d", userID, time.Now().Unix()),
			UserID:        userID,
			Title:         state.Data["title"].(string),
			Time:          state.Data["time"].(string),
			Days:          state.Data["days"].([]string),
			Note:          state.Data["note"].(string),
			ReminderType:  state.Data["reminderType"].(string),
			ReminderTimes: []int{60, 30, 5}, // Default: 1 jam, 30 menit, 5 menit sebelum
			ReminderSent:  make(map[string]bool),
		}

		if err := b.storage.AddSchedule(schedule); err != nil {
			b.sendMessage(userID, fmt.Sprintf("Error: %v", err))
		} else {
			typeStr := "Berkali-kali"
			if schedule.ReminderType == "once" {
				typeStr = "Sekali"
			}
			b.sendMessage(userID, fmt.Sprintf("✅ Jadwal berhasil ditambahkan!\n📌 %s\n⏰ Reminder: %s (1h, 30m, 5m sebelum waktu yang ditentukan)", schedule.Title, typeStr))
			b.scheduleReminder(schedule)
		}

		delete(b.userState, userID)

	case "edit_id":
		schedule, err := b.storage.GetScheduleByTitle(userID, text)
		if err != nil {
			b.sendMessage(userID, "Schedule tidak ditemukan.")
			delete(b.userState, userID)
			return
		}

		state.Data["schedule"] = schedule
		state.Action = "edit_field"
		b.userState[userID] = state
		b.sendMessageWithKeyboard(userID, "Pilih field yang ingin diubah:", getFieldKeyboard())

	case "edit_title":
		schedule, err := b.storage.GetScheduleByTitle(userID, text)
		if err != nil {
			b.sendMessage(userID, "❌ Jadwal dengan judul tersebut tidak ditemukan.")
			delete(b.userState, userID)
			return
		}

		state.Data["schedule"] = schedule
		state.Action = "edit_field"
		b.userState[userID] = state
		b.sendMessageWithKeyboard(userID, "Pilih field yang ingin diubah:", getFieldKeyboard())

	case "edit_field":
		// Use fieldMap to convert button text to field name
		fieldMap := map[string]string{
			"1": "title", "2": "time", "3": "days", "4": "note",
			"1️⃣ Title": "title", "2️⃣ Waktu": "time", "3️⃣ Hari": "days", "4️⃣ Catatan": "note",
		}
		
		field := ""
		if f, exists := fieldMap[text]; exists {
			field = f
		} else {
			b.sendReplyMessage(userID, "Field tidak valid. Pilih dari tombol yang tersedia.")
			b.sendMessageWithKeyboard(userID, "Pilih field yang ingin diubah:", getFieldKeyboard())
			return
		}

		state.Data["field"] = field
		state.Action = "edit_value"
		b.userState[userID] = state
		
		switch field {
		case "title":
			b.sendMessageWithKeyboard(userID, fmt.Sprintf("Masukkan nilai baru untuk %s:", field), getSkipKeyboard())
		case "time":
			b.sendMessageWithKeyboard(userID, fmt.Sprintf("Pilih nilai baru untuk %s:", field), getTimeKeyboard())
		case "days":
			b.sendMessageWithKeyboard(userID, fmt.Sprintf("Pilih nilai baru untuk %s:", field), getDaysKeyboard())
		case "note":
			b.sendMessageWithKeyboard(userID, fmt.Sprintf("Masukkan catatan:", field), getNoteKeyboard())
		}

	case "edit_value":
		schedule := state.Data["schedule"].(*storage.Schedule)
		field := state.Data["field"].(string)

		switch field {
		case "title":
			// Check if new title already exists (but allow same title)
			if text != schedule.Title && b.storage.IsTitleExists(userID, text) {
				b.sendReplyMessage(userID, "❌ Judul sudah digunakan. Gunakan judul yang berbeda.")
				b.sendMessageWithKeyboard(userID, "Masukkan judul baru:", getSkipKeyboard())
				return
			}
			schedule.Title = text
		case "time":
			if !isValidTime(text) {
				b.sendReplyMessage(userID, "Format waktu tidak valid.")
				return
			}
			schedule.Time = text
		case "days":
			days := parsedays(text)
			if len(days) == 0 {
				b.sendReplyMessage(userID, "Format hari tidak valid.")
				return
			}
			schedule.Days = days
		case "note":
			if text == "-" {
				schedule.Note = ""
			} else {
				schedule.Note = text
			}
		}

		if err := b.storage.UpdateSchedule(schedule); err != nil {
			b.sendMessage(userID, fmt.Sprintf("Error: %v", err))
			delete(b.userState, userID)
		} else {
			b.sendMessage(userID, "✅ "+field+" berhasil diperbarui!")
			state.Action = "edit_continue"
			b.userState[userID] = state
			b.sendMessageWithKeyboard(userID, "Ingin melanjutkan edit field lain?", getEditContinueKeyboard())
		}

	case "edit_continue":
		if text == "✏️ Lanjut Edit" {
			state.Action = "edit_field"
			b.userState[userID] = state
			b.sendMessageWithKeyboard(userID, "Pilih field yang ingin diubah:", getFieldKeyboard())
		} else if text == "✅ Selesai" {
			b.sendMessage(userID, "Perubahan jadwal selesai. Ketik /help untuk bantuan.")
			delete(b.userState, userID)
		} else {
			b.sendReplyMessage(userID, "Pilihan tidak valid. Pilih dari tombol yang tersedia.")
			b.sendMessageWithKeyboard(userID, "Ingin melanjutkan edit field lain?", getEditContinueKeyboard())
		}

	case "delete_id":
		if err := b.storage.DeleteSchedule(text); err != nil {
			b.sendMessage(userID, "Schedule tidak ditemukan.")
		} else {
			b.sendMessage(userID, "✅ Jadwal berhasil dihapus!")
		}
		delete(b.userState, userID)

	case "delete_title":
		schedule, err := b.storage.GetScheduleByTitle(userID, text)
		if err != nil {
			b.sendMessage(userID, "❌ Jadwal dengan judul tersebut tidak ditemukan.")
			delete(b.userState, userID)
			return
		}

		if err := b.storage.DeleteSchedule(schedule.ID); err != nil {
			b.sendMessage(userID, "Gagal menghapus jadwal.")
		} else {
			b.sendMessage(userID, "✅ Jadwal berhasil dihapus!")
		}
		delete(b.userState, userID)
	}
}

func (b *Bot) listSchedules(userID int64) {
	schedules := b.storage.GetUserSchedules(userID)

	if len(schedules) == 0 {
		b.sendMessage(userID, "Anda belum memiliki jadwal. Gunakan /add untuk membuat jadwal baru.")
		return
	}

	var text strings.Builder
	text.WriteString("📅 Jadwal Anda:\n\n")

	for _, s := range schedules {
		text.WriteString(fmt.Sprintf("📌 Judul: %s\n", s.Title))
		text.WriteString(fmt.Sprintf("⏰ Waktu: %s\n", s.Time))
		text.WriteString(fmt.Sprintf("📆 Hari: %s\n", strings.Join(s.Days, ", ")))
		if s.Note != "" {
			text.WriteString(fmt.Sprintf("📝 Catatan: %s\n", s.Note))
		}
		text.WriteString("\n")
	}

	b.sendMessageHTML(userID, text.String())
}

func (b *Bot) scheduleReminder(schedule *storage.Schedule) {
	for _, day := range schedule.Days {
		// Parse hour and minute from time string
		var scheduleHour, scheduleMin int
		fmt.Sscanf(schedule.Time, "%d:%d", &scheduleHour, &scheduleMin)
		
		// Schedule reminders for each reminder time
		for _, reminderMinutes := range schedule.ReminderTimes {
			// Calculate reminder time
			reminderHour := scheduleHour
			reminderMin := scheduleMin - reminderMinutes
			
			if reminderMin < 0 {
				reminderMin += 60
				reminderHour--
				if reminderHour < 0 {
					reminderHour = 23
				}
			}
			
			weekday := dayToCronDay(day)
			cronExpression := fmt.Sprintf("%d %d * * %s", reminderMin, reminderHour, weekday)
			
			scheduleID := schedule.ID
			reminderKey := fmt.Sprintf("%s_%dm", scheduleID, reminderMinutes)
			
			_, err := b.cron.AddFunc(cronExpression, func() {
				// Refresh schedule dari storage untuk get latest data
				latestSchedule, err := b.storage.GetSchedule(scheduleID)
				if err != nil {
					return
				}
				
				note := latestSchedule.Note
				if note == "" {
					note = "(tanpa catatan)"
				}
				
				// Check if reminder already sent (for "once" type)
				if latestSchedule.ReminderType == "once" {
					if latestSchedule.ReminderSent[reminderKey] {
						return // Reminder sudah pernah dikirim
					}
				}
				
				// Send reminder
				reminderText := fmt.Sprintf("⏰ Pengingat %d menit sebelum:\n📌 %s\n%s", reminderMinutes, latestSchedule.Title, note)
				b.sendMessage(latestSchedule.UserID, reminderText)
				
				// Mark as sent if type is "once"
				if latestSchedule.ReminderType == "once" {
					if latestSchedule.ReminderSent == nil {
						latestSchedule.ReminderSent = make(map[string]bool)
					}
					latestSchedule.ReminderSent[reminderKey] = true
					b.storage.UpdateSchedule(latestSchedule)
					
					// If all reminders sent, delete the schedule
					if len(latestSchedule.ReminderSent) == len(latestSchedule.ReminderTimes) {
						b.storage.DeleteSchedule(scheduleID)
					}
				}
			})
			
			if err != nil {
				log.Printf("Error scheduling reminder: %v\n", err)
			}
		}
	}
}

func (b *Bot) sendMessage(userID int64, text string) {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	b.api.Send(msg)
}

func (b *Bot) sendReplyMessage(userID int64, text string) {
	msg := tgbotapi.NewMessage(userID, text)
	b.api.Send(msg)
}

func (b *Bot) sendMessageHTML(userID int64, text string) {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	b.api.Send(msg)
}

func (b *Bot) sendMessageWithKeyboard(userID int64, text string, keyboard tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

func (b *Bot) Stop() {
	b.cron.Stop()
}

// Helper functions

func getHelpText() string {
	return `🤖 Schedule Bot - Bantuan

Perintah yang tersedia:
/add - Tambah jadwal baru
/list - Lihat semua jadwal
/edit - Ubah jadwal
/delete - Hapus jadwal
/help - Tampilkan bantuan ini

Contoh penggunaan:
1. Ketik /add
2. Ikuti petunjuk untuk membuat jadwal baru
3. Bot akan mengingatkan Anda sesuai jadwal

Untuk pertanyaan, silakan hubungi @FtrRahman`
}

func isValidTime(t string) bool {
	parts := strings.Split(t, ":")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) == 2 && len(parts[1]) == 2
}

func parsedays(text string) []string {
	days := strings.Split(text, ",")
	validDays := map[string]bool{
		"Monday": true, "Tuesday": true, "Wednesday": true, "Thursday": true,
		"Friday": true, "Saturday": true, "Sunday": true,
	}

	// Map keyboard button format ke format valid
	dayButtonMap := map[string]string{
		"Senin (Monday)":       "Monday",
		"Selasa (Tuesday)":     "Tuesday",
		"Rabu (Wednesday)":     "Wednesday",
		"Kamis (Thursday)":     "Thursday",
		"Jumat (Friday)":       "Friday",
		"Sabtu (Saturday)":     "Saturday",
		"Minggu (Sunday)":      "Sunday",
	}

	var result []string
	for _, day := range days {
		day = strings.TrimSpace(day)
		
		// Try direct match first
		if validDays[day] {
			result = append(result, day)
			continue
		}

		// Try button format mapping
		if mappedDay, exists := dayButtonMap[day]; exists {
			result = append(result, mappedDay)
			continue
		}
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func convertTimeToMinHour(timeStr string) string {
	// Convert "09:30" to "30 9"
	parts := strings.Split(timeStr, ":")
	return fmt.Sprintf("%s %s", parts[1], parts[0])
}

func dayToCronDay(day string) string {
	dayMap := map[string]int{
		"Sunday":    0,
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
	}
	return fmt.Sprintf("%d", dayMap[day])
}

// Keyboard helper functions

func getTimeKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("00:00"),
			tgbotapi.NewKeyboardButton("01:00"),
			tgbotapi.NewKeyboardButton("02:00"),
			tgbotapi.NewKeyboardButton("03:00"),
			tgbotapi.NewKeyboardButton("04:00"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("05:00"),
			tgbotapi.NewKeyboardButton("06:00"),
			tgbotapi.NewKeyboardButton("07:00"),
			tgbotapi.NewKeyboardButton("08:00"),
			tgbotapi.NewKeyboardButton("09:00"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("10:00"),
			tgbotapi.NewKeyboardButton("11:00"),
			tgbotapi.NewKeyboardButton("12:00"),
			tgbotapi.NewKeyboardButton("13:00"),
			tgbotapi.NewKeyboardButton("14:00"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("15:00"),
			tgbotapi.NewKeyboardButton("16:00"),
			tgbotapi.NewKeyboardButton("17:00"),
			tgbotapi.NewKeyboardButton("18:00"),
			tgbotapi.NewKeyboardButton("19:00"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("20:00"),
			tgbotapi.NewKeyboardButton("21:00"),
			tgbotapi.NewKeyboardButton("22:00"),
			tgbotapi.NewKeyboardButton("23:00"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❌ Batal"),
		),
	)
}

func getDaysKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Senin (Monday)"),
			tgbotapi.NewKeyboardButton("Selasa (Tuesday)"),
			tgbotapi.NewKeyboardButton("Rabu (Wednesday)"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Kamis (Thursday)"),
			tgbotapi.NewKeyboardButton("Jumat (Friday)"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Sabtu (Saturday)"),
			tgbotapi.NewKeyboardButton("Minggu (Sunday)"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔄 Selesai Pilih"),
			tgbotapi.NewKeyboardButton("❌ Batal"),
		),
	)
}

func getNoteKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Tidak ada catatan"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❌ Batal"),
		),
	)
}

func getSkipKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❌ Batal"),
		),
	)
}

func getFieldKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("1️⃣ Title"),
			tgbotapi.NewKeyboardButton("2️⃣ Waktu"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("3️⃣ Hari"),
			tgbotapi.NewKeyboardButton("4️⃣ Catatan"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❌ Batal"),
		),
	)
}

func getReminderTypeKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔔 Sekali"),
			tgbotapi.NewKeyboardButton("🔊 Berkali-kali"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❌ Batal"),
		),
	)
}

func getEditContinueKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✏️ Lanjut Edit"),
			tgbotapi.NewKeyboardButton("✅ Selesai"),
		),
	)
}
