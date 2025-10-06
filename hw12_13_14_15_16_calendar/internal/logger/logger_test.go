package logger

import (
	"bytes"
	"strings"
	"testing"

	internalcfg "github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_calendar/internal/config"
)

func Test_parseLevel(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"debug", "DEBUG"},
		{"info", "INFO"},
		{"", "INFO"},
		{"warn", "WARN"},
		{"warning", "WARN"},
		{"error", "ERROR"},
		{"unknown", "INFO"},
	}
	for _, tc := range cases {
		lvl := parseLevel(tc.in)
		if lvl.String() != tc.want {
			t.Fatalf("parseLevel(%q) = %s, want %s", tc.in, lvl, tc.want)
		}
	}
}

func Test_Info_Error_WriteAndFiltering(t *testing.T) {
	var buf bytes.Buffer
	cfg := &internalcfg.Logger{Level: "info"}
	l := NewWithWriter(&buf, cfg)

	l.Info("test info")
	if got := buf.String(); !strings.Contains(got, "test info") || !strings.Contains(got, "level=INFO") {
		t.Fatalf("Info not written correctly, got: %s", got)
	}

	buf.Reset()
	l.Error("error test")
	if got := buf.String(); !strings.Contains(got, "error test") || !strings.Contains(got, "level=ERROR") {
		t.Fatalf("Error not written correctly, got: %s", got)
	}

	// Filtering: set level=ERROR, Info should not appear
	buf.Reset()
	cfg.Level = "error"
	l = NewWithWriter(&buf, cfg)
	l.Info("should be filtered")
	if got := buf.String(); got != "" {
		t.Fatalf("Info should be filtered at error level, got: %q", got)
	}

	l.Error("error visible")
	if got := buf.String(); !strings.Contains(got, "error visible") {
		t.Fatalf("Error should be logged at error level, got: %q", got)
	}
}
