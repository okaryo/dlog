package render

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/okaryo/dlog/internal/model"
)

type formattedLog struct {
	timestamp time.Time
	text      string
}

func Text(dayLog *model.DayLog) (string, error) {
	logs, err := parseLogs(dayLog)
	if err != nil {
		return "", err
	}

	slices.SortFunc(logs, func(a, b formattedLog) int {
		if a.timestamp.Equal(b.timestamp) {
			return 0
		}
		if a.timestamp.After(b.timestamp) {
			return -1
		}
		return 1
	})

	if len(logs) == 0 {
		return dayLog.Date, nil
	}

	lines := []string{dayLog.Date, ""}
	for _, entry := range logs {
		lines = append(lines, fmt.Sprintf("%s %s", entry.timestamp.In(time.Local).Format("15:04"), entry.text))
	}

	return strings.Join(lines, "\n"), nil
}

func Markdown(dayLog *model.DayLog) (string, error) {
	logs, err := parseLogs(dayLog)
	if err != nil {
		return "", err
	}

	slices.SortFunc(logs, func(a, b formattedLog) int {
		if a.timestamp.Equal(b.timestamp) {
			return 0
		}
		if a.timestamp.Before(b.timestamp) {
			return -1
		}
		return 1
	})

	lines := []string{"# " + dayLog.Date}
	for _, entry := range logs {
		lines = append(lines, fmt.Sprintf("- %s %s", entry.timestamp.In(time.Local).Format("15:04"), entry.text))
	}

	return strings.Join(lines, "\n"), nil
}

func parseLogs(dayLog *model.DayLog) ([]formattedLog, error) {
	logs := make([]formattedLog, 0, len(dayLog.Logs))
	for _, entry := range dayLog.Logs {
		parsed, err := time.Parse(time.RFC3339, entry.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("parse timestamp %q: %w", entry.Timestamp, err)
		}

		logs = append(logs, formattedLog{
			timestamp: parsed,
			text:      entry.Text,
		})
	}

	return logs, nil
}
