package configuration

import (
	"encoding/json"
	"fmt"
	"io"
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
	Writer io.Writer
}

func (li *LogInterceptor) Write(p []byte) (n int, err error) {
	var entry LogEntry
	err = json.Unmarshal(p, &entry)
	if err != nil {
		return li.Writer.Write(p)
	}

	t, err := time.Parse(time.RFC3339Nano, entry.Time)
	if err != nil {
		t = time.Now()
	}

	formatted := fmt.Sprintf("%s %s %s\n", t.Format("2006/01/02 15:04:05"), entry.Level, entry.Message)
	return li.Writer.Write([]byte(formatted))
}
