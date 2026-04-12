package model

type DayLog struct {
	Date string     `json:"date"`
	Logs []LogEntry `json:"logs"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Text      string `json:"text"`
}
