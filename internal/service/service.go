package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/okaryo/dlog/internal/model"
	"github.com/okaryo/dlog/internal/storage"
)

const dayLayout = "2006-01-02"

type Service struct {
	store *storage.Store
	now   func() time.Time
}

func New(store *storage.Store) *Service {
	return &Service{
		store: store,
		now:   time.Now,
	}
}

func NewWithNow(store *storage.Store, now func() time.Time) *Service {
	return &Service{
		store: store,
		now:   now,
	}
}

func (s *Service) AddTodayLog(text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return fmt.Errorf("log text cannot be empty")
	}

	now := s.now()
	dayLog, err := s.store.LoadDay(now)
	if err != nil {
		return err
	}

	expectedDate := now.Format(dayLayout)
	if dayLog.Date != expectedDate {
		return fmt.Errorf("day log date mismatch: expected %s, got %s", expectedDate, dayLog.Date)
	}

	dayLog.Logs = append(dayLog.Logs, model.LogEntry{
		Timestamp: now.Format(time.RFC3339),
		Text:      text,
	})

	return s.store.SaveDay(dayLog)
}

func (s *Service) AmendTodayLog(text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return fmt.Errorf("log text cannot be empty")
	}

	now := s.now()
	dayLog, err := s.store.LoadDay(now)
	if err != nil {
		return err
	}

	expectedDate := now.Format(dayLayout)
	if dayLog.Date != expectedDate {
		return fmt.Errorf("day log date mismatch: expected %s, got %s", expectedDate, dayLog.Date)
	}

	if len(dayLog.Logs) == 0 {
		return fmt.Errorf("no logs to amend for today")
	}

	dayLog.Logs[len(dayLog.Logs)-1].Text = text

	return s.store.SaveDay(dayLog)
}

func (s *Service) GetTodayLog() (*model.DayLog, error) {
	return s.getLogByTime(s.now())
}

func (s *Service) GetLogByDate(date string) (*model.DayLog, error) {
	parsedDate, err := s.parseDate(date)
	if err != nil {
		return nil, err
	}

	return s.getLogByTime(parsedDate)
}

func (s *Service) parseDate(date string) (time.Time, error) {
	switch strings.ToLower(strings.TrimSpace(date)) {
	case "today":
		return s.now(), nil
	case "yesterday":
		return s.now().AddDate(0, 0, -1), nil
	}

	parsedDate, err := time.ParseInLocation(dayLayout, date, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("date must be YYYY-MM-DD, today, or yesterday: %q", date)
	}

	return parsedDate, nil
}

func (s *Service) getLogByTime(target time.Time) (*model.DayLog, error) {
	dayLog, err := s.store.LoadDay(target)
	if err != nil {
		return nil, err
	}

	expectedDate := target.Format(dayLayout)
	if dayLog.Date != expectedDate {
		return nil, fmt.Errorf("day log date mismatch: expected %s, got %s", expectedDate, dayLog.Date)
	}

	return dayLog, nil
}
