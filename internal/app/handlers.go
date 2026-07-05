package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloudnativeapp/internal/store"
)

func (a *App) handleUI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_ = a.store.Write(store.Event{
		Type:      "ui_view",
		Hostname:  a.identity.Hostname(),
		IPAddress: a.identity.IPAddress(),
		Pod:       os.Getenv("HOSTNAME"),
		Path:      r.URL.Path,
		Method:    r.Method,
		Remote:    r.RemoteAddr,
		Note:      "dashboard rendered",
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := dashboardData{
		Hostname:  a.identity.Hostname(),
		IPAddress: a.identity.IPAddress(),
		Pod:       os.Getenv("HOSTNAME"),
		Uptime:    time.Since(a.startedAt).Round(time.Second).String(),
	}

	if err := a.dashboard.Render(w, data); err != nil {
		http.Error(w, fmt.Sprintf("unable to render dashboard: %v", err), http.StatusInternalServerError)
	}
}

func (a *App) handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_ = a.store.Write(store.Event{
		Type:      "api_info_view",
		Hostname:  a.identity.Hostname(),
		IPAddress: a.identity.IPAddress(),
		Pod:       os.Getenv("HOSTNAME"),
		Path:      r.URL.Path,
		Method:    r.Method,
		Remote:    r.RemoteAddr,
		Note:      "machine info requested",
	})

	writeJSON(w, map[string]string{
		"status":     "success",
		"message":    "Go cloud-native app is running",
		"hostname":   a.identity.Hostname(),
		"ip_address": a.identity.IPAddress(),
		"pod_name":   os.Getenv("HOSTNAME"),
		"timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
		"uptime":     time.Since(a.startedAt).Round(time.Second).String(),
	})
}

func (a *App) handlePod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	podName := os.Getenv("HOSTNAME")
	_ = a.store.Write(store.Event{
		Type:      "pod_view",
		Hostname:  a.identity.Hostname(),
		IPAddress: a.identity.IPAddress(),
		Pod:       podName,
		Path:      r.URL.Path,
		Method:    r.Method,
		Remote:    r.RemoteAddr,
		Note:      "pod identity requested",
	})

	writeJSON(w, map[string]string{
		"status":     "success",
		"message":    "request served by this pod",
		"hostname":   a.identity.Hostname(),
		"ip_address": a.identity.IPAddress(),
		"pod_name":   podName,
		"timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
	})
}

func (a *App) handlePanic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_ = a.store.Write(store.Event{
		Type:      "panic_triggered",
		Hostname:  a.identity.Hostname(),
		IPAddress: a.identity.IPAddress(),
		Pod:       os.Getenv("HOSTNAME"),
		Path:      r.URL.Path,
		Method:    r.Method,
		Remote:    r.RemoteAddr,
		Note:      "dashboard requested process crash",
	})

	writeJSON(w, map[string]string{
		"status":     "accepted",
		"message":    "panic scheduled; the process will exit",
		"hostname":   a.identity.Hostname(),
		"ip_address": a.identity.IPAddress(),
		"pod_name":   os.Getenv("HOSTNAME"),
		"timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
	})

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	go func() {
		time.Sleep(150 * time.Millisecond)
		panic("simulated crash from dashboard button")
	}()
}

func (a *App) handleWork(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	seconds := 15 * time.Second
	started := time.Now().UTC()
	deadline := time.Now().Add(seconds)
	iterations := 0
	var checksum uint64

	_ = a.store.Write(store.Event{
		Type:      "work_started",
		Hostname:  a.identity.Hostname(),
		IPAddress: a.identity.IPAddress(),
		Pod:       os.Getenv("HOSTNAME"),
		Path:      r.URL.Path,
		Method:    r.Method,
		Remote:    r.RemoteAddr,
		Note:      "cpu burn started to encourage hpa scale up",
	})

	for time.Now().Before(deadline) {
		for i := 0; i < 750000; i++ {
			checksum = checksum*1664525 + uint64(i) + 1013904223
		}
		iterations++
	}

	_ = a.store.Write(store.Event{
		Type:      "work_finished",
		Hostname:  a.identity.Hostname(),
		IPAddress: a.identity.IPAddress(),
		Pod:       os.Getenv("HOSTNAME"),
		Path:      r.URL.Path,
		Method:    r.Method,
		Remote:    r.RemoteAddr,
		Note:      "cpu burn completed",
	})

	writeJSON(w, map[string]string{
		"status":      "success",
		"message":     "cpu burn completed; hit again to help hpa move to the second pod",
		"hostname":    a.identity.Hostname(),
		"ip_address":  a.identity.IPAddress(),
		"pod_name":    os.Getenv("HOSTNAME"),
		"started_at":  started.Format(time.RFC3339Nano),
		"finished_at": time.Now().UTC().Format(time.RFC3339Nano),
		"duration":    time.Since(started).Round(time.Millisecond).String(),
		"iterations":  fmt.Sprintf("%d", iterations),
		"checksum":    fmt.Sprintf("%d", checksum),
	})
}

func (a *App) handleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	events, err := a.store.ReadAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to read state: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]any{
		"status":      "success",
		"log_path":    a.store.Path(),
		"event_count": len(events),
		"events":      events,
	})
}

func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(v)
}
