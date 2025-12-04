package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Schedule struct {
	ID              string    `json:"id"`
	UserID          int64     `json:"user_id"`
	Title           string    `json:"title"`
	Time            string    `json:"time"` // Format: HH:MM
	Days            []string  `json:"days"` // [Monday, Tuesday, ...]
	Note            string    `json:"note"`
	ReminderType    string    `json:"reminder_type"`    // "once" atau "recurring"
	ReminderTimes   []int     `json:"reminder_times"`   // [60, 30, 5] dalam menit sebelum jadwal
	ReminderSent    map[string]bool `json:"reminder_sent"` // Track reminder yang sudah dikirim untuk tipe "once"
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserSchedules struct {
	Schedules map[string]*Schedule `json:"schedules"`
	mu        sync.RWMutex
	filePath  string
}

func NewUserSchedules(filePath string) (*UserSchedules, error) {
	us := &UserSchedules{
		Schedules: make(map[string]*Schedule),
		filePath:  filePath,
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("gagal membuat direktori: %w", err)
	}

	// Load existing data
	if err := us.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return us, nil
}

func (us *UserSchedules) load() error {
	us.mu.Lock()
	defer us.mu.Unlock()

	data, err := os.ReadFile(us.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	var schedules map[string]*Schedule
	if err := json.Unmarshal(data, &schedules); err != nil {
		return fmt.Errorf("gagal parse JSON: %w", err)
	}

	us.Schedules = schedules
	return nil
}

func (us *UserSchedules) Save() error {
	us.mu.RLock()
	defer us.mu.RUnlock()

	data, err := json.MarshalIndent(us.Schedules, "", "  ")
	if err != nil {
		return fmt.Errorf("gagal marshal JSON: %w", err)
	}

	if err := os.WriteFile(us.filePath, data, 0644); err != nil {
		return fmt.Errorf("gagal menyimpan file: %w", err)
	}

	return nil
}

func (us *UserSchedules) AddSchedule(schedule *Schedule) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()
	us.Schedules[schedule.ID] = schedule

	return us.saveUnlocked()
}

func (us *UserSchedules) UpdateSchedule(schedule *Schedule) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	if _, exists := us.Schedules[schedule.ID]; !exists {
		return fmt.Errorf("schedule tidak ditemukan")
	}

	schedule.UpdatedAt = time.Now()
	us.Schedules[schedule.ID] = schedule

	return us.saveUnlocked()
}

func (us *UserSchedules) DeleteSchedule(id string) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	if _, exists := us.Schedules[id]; !exists {
		return fmt.Errorf("schedule tidak ditemukan")
	}

	delete(us.Schedules, id)
	return us.saveUnlocked()
}

func (us *UserSchedules) GetUserSchedules(userID int64) []*Schedule {
	us.mu.RLock()
	defer us.mu.RUnlock()

	var result []*Schedule
	for _, schedule := range us.Schedules {
		if schedule.UserID == userID {
			result = append(result, schedule)
		}
	}

	return result
}

func (us *UserSchedules) GetSchedule(id string) (*Schedule, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	schedule, exists := us.Schedules[id]
	if !exists {
		return nil, fmt.Errorf("schedule tidak ditemukan")
	}

	return schedule, nil
}

func (us *UserSchedules) GetScheduleByTitle(userID int64, title string) (*Schedule, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	for _, schedule := range us.Schedules {
		if schedule.UserID == userID && schedule.Title == title {
			return schedule, nil
		}
	}

	return nil, fmt.Errorf("schedule dengan judul '%s' tidak ditemukan", title)
}

func (us *UserSchedules) IsTitleExists(userID int64, title string) bool {
	us.mu.RLock()
	defer us.mu.RUnlock()

	for _, schedule := range us.Schedules {
		if schedule.UserID == userID && schedule.Title == title {
			return true
		}
	}

	return false
}

func (us *UserSchedules) saveUnlocked() error {
	data, err := json.MarshalIndent(us.Schedules, "", "  ")
	if err != nil {
		return fmt.Errorf("gagal marshal JSON: %w", err)
	}

	if err := os.WriteFile(us.filePath, data, 0644); err != nil {
		return fmt.Errorf("gagal menyimpan file: %w", err)
	}

	return nil
}
