package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/okaryo/dlog/internal/model"
)

const dayLayout = "2006-01-02"

type Store struct {
	baseDir string
}

func NewDefaultStore() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolve home directory: %w", err)
	}

	return NewStore(filepath.Join(homeDir, ".dlog")), nil
}

func NewStore(baseDir string) *Store {
	return &Store{baseDir: baseDir}
}

func (s *Store) BaseDir() string {
	return s.baseDir
}

func (s *Store) LoadDay(date time.Time) (*model.DayLog, error) {
	path := s.DayPath(date)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &model.DayLog{
				Date: date.In(time.Local).Format(dayLayout),
				Logs: []model.LogEntry{},
			}, nil
		}

		return nil, fmt.Errorf("read day log file %q: %w", path, err)
	}

	var dayLog model.DayLog
	if err := json.Unmarshal(data, &dayLog); err != nil {
		return nil, fmt.Errorf("parse day log file %q: %w", path, err)
	}

	if dayLog.Date == "" {
		return nil, fmt.Errorf("day log file %q is missing date", path)
	}

	return &dayLog, nil
}

func (s *Store) SaveDay(dayLog *model.DayLog) error {
	if err := os.MkdirAll(s.baseDir, 0o755); err != nil {
		return fmt.Errorf("create storage directory %q: %w", s.baseDir, err)
	}

	path := filepath.Join(s.baseDir, dayLog.Date+".json")
	data, err := json.MarshalIndent(dayLog, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal day log: %w", err)
	}
	data = append(data, '\n')

	tempFile, err := os.CreateTemp(s.baseDir, filepath.Base(path)+".tmp.*")
	if err != nil {
		return fmt.Errorf("create temp file for %q: %w", path, err)
	}

	tempName := tempFile.Name()
	cleanup := func() {
		_ = os.Remove(tempName)
	}

	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		cleanup()
		return fmt.Errorf("write temp file %q: %w", tempName, err)
	}

	if err := tempFile.Close(); err != nil {
		cleanup()
		return fmt.Errorf("close temp file %q: %w", tempName, err)
	}

	if err := os.Rename(tempName, path); err != nil {
		cleanup()
		return fmt.Errorf("replace day log file %q: %w", path, err)
	}

	return nil
}

func (s *Store) DayPath(date time.Time) string {
	return filepath.Join(s.baseDir, date.In(time.Local).Format(dayLayout)+".json")
}
