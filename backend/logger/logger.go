package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
)

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiRed    = "\033[31m"
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
	ansiGray   = "\033[90m"
)

type colorHandler struct {
	mu  sync.Mutex
	out io.Writer
	lvl slog.Level
}

func New(level slog.Level) slog.Handler {
	return &colorHandler{out: os.Stdout, lvl: level}
}

func Init() {
	slog.SetDefault(slog.New(New(slog.LevelInfo)))
}

func (h *colorHandler) Enabled(_ context.Context, l slog.Level) bool { return l >= h.lvl }
func (h *colorHandler) WithGroup(name string) slog.Handler             { return h }
func (h *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler       { return h }

func (h *colorHandler) Handle(_ context.Context, r slog.Record) error {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%s%s%s  ", ansiGray, r.Time.Format("15:04:05"), ansiReset)

	switch r.Level {
	case slog.LevelInfo:
		fmt.Fprintf(&buf, "%s%sINFO%s  ", ansiCyan, ansiBold, ansiReset)
	case slog.LevelWarn:
		fmt.Fprintf(&buf, "%s%sWARN%s  ", ansiYellow, ansiBold, ansiReset)
	case slog.LevelError:
		fmt.Fprintf(&buf, "%s%sERRO%s  ", ansiRed, ansiBold, ansiReset)
	default:
		fmt.Fprintf(&buf, "%sDBUG%s  ", ansiGray, ansiReset)
	}

	fmt.Fprintf(&buf, "%s%s%s", ansiBold, r.Message, ansiReset)

	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(&buf, "  %s%s%s=%v", ansiGray, a.Key, ansiReset, a.Value)
		return true
	})

	buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf.Bytes())
	return err
}
