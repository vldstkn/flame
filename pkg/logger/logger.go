package logger

import (
	"context"
	"encoding/json"
	"github.com/fatih/color"
	"io"
	"log"
	"log/slog"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	b, err := json.MarshalIndent(fields, "", " ")

	if err != nil {
		return err
	}
	timeStr := color.HiWhiteString(r.Time.Format("[15:05:05.000]"))

	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(string(string(b))))

	return nil
}

func NewLogger(out io.Writer) *slog.Logger {

	opt := &PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opt.SlogOpts),
		l:       log.New(out, "", 0),
	}

	logger := slog.New(h)
	return logger
}
