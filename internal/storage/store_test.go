package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/okaryo/dlog/internal/model"
)

func TestLoadDayReturnsEmptyDayWhenFileDoesNotExist(t *testing.T) {
	store := NewStore(t.TempDir())
	date := time.Date(2026, 4, 12, 0, 0, 0, 0, time.FixedZone("JST", 9*60*60))

	dayLog, err := store.LoadDay(date)
	if err != nil {
		t.Fatalf("load missing day: %v", err)
	}

	if dayLog.Date != "2026-04-12" {
		t.Fatalf("unexpected date: %s", dayLog.Date)
	}
	if len(dayLog.Logs) != 0 {
		t.Fatalf("expected empty logs, got %d", len(dayLog.Logs))
	}
}

func TestSaveDayWritesIndentedJSONAtomically(t *testing.T) {
	store := NewStore(t.TempDir())
	dayLog := &model.DayLog{
		Date: "2026-04-12",
		Logs: []model.LogEntry{
			{Timestamp: "2026-04-12T10:03:21+09:00", Text: "task progress update"},
		},
	}

	if err := store.SaveDay(dayLog); err != nil {
		t.Fatalf("save day: %v", err)
	}

	entries, err := os.ReadDir(store.BaseDir())
	if err != nil {
		t.Fatalf("read storage dir: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected one final file, got %d", len(entries))
	}

	if entries[0].Name() != "2026-04-12.json" {
		t.Fatalf("unexpected file name: %s", entries[0].Name())
	}

	path := filepath.Join(store.BaseDir(), entries[0].Name())
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read saved file: %v", err)
	}

	want := "{\n  \"date\": \"2026-04-12\",\n  \"logs\": [\n    {\n      \"timestamp\": \"2026-04-12T10:03:21+09:00\",\n      \"text\": \"task progress update\"\n    }\n  ]\n}\n"
	if string(data) != want {
		t.Fatalf("unexpected file content:\nwant:\n%s\ngot:\n%s", want, string(data))
	}
}
