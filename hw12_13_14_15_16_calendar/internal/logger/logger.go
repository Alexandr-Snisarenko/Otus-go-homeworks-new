package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	internalcfg "github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_calendar/internal/config"
)

type Logger struct {
	innerLog *slog.Logger
}

// Non-context wrappers.
func (l *Logger) Info(msg string, v ...any) {
	l.innerLog.Info(msg, v...)
}

func (l *Logger) Error(msg string, v ...any) {
	l.innerLog.Error(msg, v...)
}

func (l *Logger) Fatal(msg string, v ...any) {
	l.innerLog.Error(msg, v...)
	os.Exit(1)
}

func (l *Logger) Debug(msg string, v ...any) {
	l.innerLog.Debug(msg, v...)
}

func (l *Logger) Warn(msg string, v ...any) {
	l.innerLog.Warn(msg, v...)
}

// Context-aware wrappers.
func (l *Logger) InfoCtx(ctx context.Context, msg string, v ...any) {
	l.innerLog.InfoContext(ctx, msg, v...)
}

func (l *Logger) ErrorCtx(ctx context.Context, msg string, v ...any) {
	l.innerLog.ErrorContext(ctx, msg, v...)
}

func (l *Logger) DebugCtx(ctx context.Context, msg string, v ...any) {
	l.innerLog.DebugContext(ctx, msg, v...)
}

func (l *Logger) WarnCtx(ctx context.Context, msg string, v ...any) {
	l.innerLog.WarnContext(ctx, msg, v...)
}

// New - создаём логгер на базе slog
// По умолчанию level = info
// Вывод лога - в stdout.
func New(logCfg *internalcfg.Logger) *Logger {
	lvl := parseLevel(logCfg.Level)
	hdlr := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	return &Logger{innerLog: slog.New(hdlr)}
}

// NewWithWriter позволяет писать в io.Writer. Для тестирования
// Если writer не задан - пишем в os.Stdout.
func NewWithWriter(w io.Writer, logCfg *internalcfg.Logger) *Logger {
	if w == nil {
		w = os.Stdout
	}
	lvl := parseLevel(logCfg.Level)
	hdlr := slog.NewTextHandler(w, &slog.HandlerOptions{Level: lvl})
	return &Logger{innerLog: slog.New(hdlr)}
}

func parseLevel(s string) slog.Level {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "debug":
		return slog.LevelDebug
	case "info", "":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
