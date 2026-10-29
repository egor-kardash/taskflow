package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestNewRejectsUnknownLevel(t *testing.T) {
	if _, err := New("unknownLevel", "text"); err == nil {
		t.Fatal("expected an error for unknown level")
	}
}

func TestNewRejectsUnknownFormat(t *testing.T) {
	if _, err := New("info", "xml"); err == nil {
		t.Fatal("expected an error for unknown format")
	}
}

func TestNewHandlerWritesJSON(t *testing.T) {
	var buf bytes.Buffer

	handler, err := newHandler(&buf, "debug", "json")
	if err != nil {
		t.Fatalf("newHandler: %v", err)
	}

	slog.New(handler).Info("hello", "key", "value")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal log line: %v", err)
	}
	if got["msg"] != "hello" {
		t.Fatalf("msg = %v, want hello", got["msg"])
	}
	if got["key"] != "value" {
		t.Fatalf("key = %v, want value", got["key"])
	}
}

func TestNewHandlerHonorsLevel(t *testing.T) {
	var buf bytes.Buffer

	handler, err := newHandler(&buf, "warn", "text")
	if err != nil {
		t.Fatalf("newHandler: %v", err)
	}

	slog.New(handler).Info("should be dropped")

	if strings.TrimSpace(buf.String()) != "" {
		t.Fatalf("expected no output below warn, got %q", buf.String())
	}
}
