package render

import (
	"testing"
	"time"

	"github.com/okaryo/dlog/internal/model"
)

func TestTextSortsLogsDescending(t *testing.T) {
	dayLog := &model.DayLog{
		Date: "2026-04-12",
		Logs: []model.LogEntry{
			{Timestamp: "2026-04-12T10:03:21+09:00", Text: "task progress update"},
			{Timestamp: "2026-04-12T11:10:00+09:00", Text: "bug fix"},
			{Timestamp: "2026-04-12T10:25:10+09:00", Text: "api design"},
		},
	}

	got, err := Text(dayLog)
	if err != nil {
		t.Fatalf("render text: %v", err)
	}

	want := "2026-04-12\n\n11:10 bug fix\n10:25 api design\n10:03 task progress update"
	if got != want {
		t.Fatalf("unexpected output:\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestMarkdownSortsLogsAscending(t *testing.T) {
	dayLog := &model.DayLog{
		Date: "2026-04-12",
		Logs: []model.LogEntry{
			{Timestamp: "2026-04-12T11:10:00+09:00", Text: "bug fix"},
			{Timestamp: "2026-04-12T10:03:21+09:00", Text: "task progress update"},
			{Timestamp: "2026-04-12T10:25:10+09:00", Text: "api design"},
		},
	}

	got, err := Markdown(dayLog)
	if err != nil {
		t.Fatalf("render markdown: %v", err)
	}

	want := "# 2026-04-12\n- 10:03 task progress update\n- 10:25 api design\n- 11:10 bug fix"
	if got != want {
		t.Fatalf("unexpected output:\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestTextUsesRecordedTimestampTimezoneForDisplay(t *testing.T) {
	originalLocal := time.Local
	time.Local = time.UTC
	t.Cleanup(func() {
		time.Local = originalLocal
	})

	dayLog := &model.DayLog{
		Date: "2026-04-12",
		Logs: []model.LogEntry{
			{Timestamp: "2026-04-12T10:03:21+09:00", Text: "task progress update"},
		},
	}

	got, err := Text(dayLog)
	if err != nil {
		t.Fatalf("render text: %v", err)
	}

	want := "2026-04-12\n\n10:03 task progress update"
	if got != want {
		t.Fatalf("unexpected output:\nwant:\n%s\ngot:\n%s", want, got)
	}
}
