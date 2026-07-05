package store

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Path string
}

type EventLog struct {
	path string
	mu   sync.Mutex
}

type Event struct {
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ip_address"`
	Pod       string `json:"pod_name,omitempty"`
	Path      string `json:"path,omitempty"`
	Method    string `json:"method,omitempty"`
	Remote    string `json:"remote,omitempty"`
	Note      string `json:"note,omitempty"`
}

func New(config Config) *EventLog {
	return &EventLog{path: config.Path}
}

func (e *EventLog) Path() string {
	return e.path
}

func (e *EventLog) Write(event Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(e.path), 0o755); err != nil {
		return err
	}

	file, err := os.OpenFile(e.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	event.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	encoded, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = file.Write(append(encoded, '\n'))
	return err
}

func (e *EventLog) ReadAll() ([]Event, error) {
	file, err := os.Open(e.path)
	if errors.Is(err, os.ErrNotExist) {
		return []Event{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.SliceStable(events, func(i, j int) bool {
		return events[i].Timestamp > events[j].Timestamp
	})

	return events, nil
}
