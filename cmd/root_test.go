package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/okaryo/dlog/internal/service"
	"github.com/okaryo/dlog/internal/storage"
	"github.com/spf13/cobra"
)

func TestRootCommandAddsLogWithSingleArgument(t *testing.T) {
	now := time.Date(2026, 4, 12, 10, 3, 21, 0, time.FixedZone("JST", 9*60*60))
	root, _, _, store := newTestRootCmd(t, func() time.Time { return now })

	root.SetArgs([]string{"  task progress update  "})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute root add: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(store.BaseDir(), "2026-04-12.json"))
	if err != nil {
		t.Fatalf("read saved file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, `"date": "2026-04-12"`) {
		t.Fatalf("saved file missing date: %s", content)
	}
	if !strings.Contains(content, `"timestamp": "2026-04-12T10:03:21+09:00"`) {
		t.Fatalf("saved file missing timestamp: %s", content)
	}
	if !strings.Contains(content, `"text": "task progress update"`) {
		t.Fatalf("saved file missing trimmed text: %s", content)
	}
}

func TestRootCommandShowsTodayLogsWithoutArguments(t *testing.T) {
	current := time.Date(2026, 4, 12, 10, 3, 0, 0, time.FixedZone("JST", 9*60*60))
	root, stdout, _, _ := newTestRootCmd(t, func() time.Time { return current })

	root.SetArgs([]string{"first task"})
	if err := root.Execute(); err != nil {
		t.Fatalf("seed first log: %v", err)
	}

	current = time.Date(2026, 4, 12, 11, 10, 0, 0, time.FixedZone("JST", 9*60*60))
	root.SetArgs([]string{"second task"})
	if err := root.Execute(); err != nil {
		t.Fatalf("seed second log: %v", err)
	}

	stdout.Reset()
	root.SetArgs(nil)
	if err := root.Execute(); err != nil {
		t.Fatalf("execute root log: %v", err)
	}

	want := "2026-04-12\n\n11:10 second task\n10:03 first task\n"
	if stdout.String() != want {
		t.Fatalf("unexpected output:\nwant:\n%s\ngot:\n%s", want, stdout.String())
	}
}

func TestLogSubcommandShowsDateOnlyForEmptyDay(t *testing.T) {
	now := time.Date(2026, 4, 12, 9, 0, 0, 0, time.FixedZone("JST", 9*60*60))
	root, stdout, _, _ := newTestRootCmd(t, func() time.Time { return now })

	root.SetArgs([]string{"log"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute log command: %v", err)
	}

	if stdout.String() != "2026-04-12\n" {
		t.Fatalf("unexpected output: %q", stdout.String())
	}
}

func TestMarkdownSubcommandShowsHeadingOnlyForEmptyDay(t *testing.T) {
	now := time.Date(2026, 4, 12, 9, 0, 0, 0, time.FixedZone("JST", 9*60*60))
	root, stdout, _, _ := newTestRootCmd(t, func() time.Time { return now })

	root.SetArgs([]string{"md"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute md command: %v", err)
	}

	if stdout.String() != "# 2026-04-12\n" {
		t.Fatalf("unexpected output: %q", stdout.String())
	}
}

func TestRootCommandRejectsEmptyText(t *testing.T) {
	now := time.Date(2026, 4, 12, 10, 3, 21, 0, time.FixedZone("JST", 9*60*60))
	root, _, _, _ := newTestRootCmd(t, func() time.Time { return now })

	root.SetArgs([]string{"   "})
	err := root.Execute()
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), "log text cannot be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func newTestRootCmd(t *testing.T, now func() time.Time) (*cobra.Command, *bytes.Buffer, *bytes.Buffer, *storage.Store) {
	t.Helper()

	baseDir := t.TempDir()
	store := storage.NewStore(baseDir)
	svc := service.NewWithNow(store, now)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	root := NewRootCmd(svc, stdout, stderr)

	return root, stdout, stderr, store
}
