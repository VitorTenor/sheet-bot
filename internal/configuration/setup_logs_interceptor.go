package configuration

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type LogEntry struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	File    string `json:"file"`
	Line    string `json:"line"`
	Message string `json:"message"`
}

type LogInterceptor struct {
	Writer    io.Writer
	AppConfig *ApplicationConfig
}

const dateFormat = "2006/01/02 15:04:05"

func formatLogEntry(entry LogEntry, groupName string) string {
	t, err := time.Parse(time.RFC3339Nano, entry.Time)
	if err != nil {
		t = time.Now()
	}

	return fmt.Sprintf("%s %s [%s][%s] %s\n", t.Format(dateFormat), entry.Level, entry.File, strings.ToUpper(groupName), entry.Message)
}

func (li *LogInterceptor) Write(p []byte) (n int, err error) {
	var entry LogEntry
	err = json.Unmarshal(p, &entry)
	if err != nil {
		return li.Writer.Write(p)
	}

	formatted := formatLogEntry(entry, li.AppConfig.WhatsApp.GroupName)
	return li.Writer.Write([]byte(formatted))
}
