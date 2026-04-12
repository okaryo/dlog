package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/okaryo/dlog/internal/storage"
)

func TestAddTodayLogCreatesAndAppendsEntries(t *testing.T) {
	now := time.Date(2026, 4, 12, 10, 3, 21, 0, time.FixedZone("JST", 9*60*60))
	store := storage.NewStore(t.TempDir())
	svc := NewWithNow(store, func() time.Time { return now })

	if err := svc.AddTodayLog("task progress update"); err != nil {
		t.Fatalf("add first log: %v", err)
	}
	if err := svc.AddTodayLog("api design"); err != nil {
		t.Fatalf("add second log: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(store.BaseDir(), "2026-04-12.json"))
	if err != nil {
		t.Fatalf("read day log file: %v", err)
	}

	content := string(data)
	if strings.Count(content, `"timestamp": "2026-04-12T10:03:21+09:00"`) != 2 {
		t.Fatalf("expected two log entries, got content: %s", content)
	}
	if !strings.Contains(content, `"text": "task progress update"`) || !strings.Contains(content, `"text": "api design"`) {
		t.Fatalf("saved file missing log texts: %s", content)
	}
}

func TestAddTodayLogFailsOnCorruptJSONWithoutChangingFile(t *testing.T) {
	now := time.Date(2026, 4, 12, 10, 3, 21, 0, time.FixedZone("JST", 9*60*60))
	store := storage.NewStore(t.TempDir())
	svc := NewWithNow(store, func() time.Time { return now })

	path := filepath.Join(store.BaseDir(), "2026-04-12.json")
	if err := os.WriteFile(path, []byte("{broken"), 0o644); err != nil {
		t.Fatalf("write corrupt file: %v", err)
	}

	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read corrupt file: %v", err)
	}

	err = svc.AddTodayLog("task progress update")
	if err == nil {
		t.Fatalf("expected error")
	}

	after, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("read after failed add: %v", readErr)
	}

	if string(before) != string(after) {
		t.Fatalf("corrupt file was modified")
	}
}
