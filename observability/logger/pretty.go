package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	AddSource   bool
	Level       slog.Leveler
	ReplaceAttr func(groups []string, attr slog.Attr) slog.Attr
}

type PrettyHandler struct {
	out   io.Writer
	opts  PrettyHandlerOptions
	attrs []slog.Attr
	group string
}

func NewPrettyHandler(out io.Writer, opts *PrettyHandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &PrettyHandlerOptions{}
	}
	return &PrettyHandler{
		out:  out,
		opts: *opts,
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	if h.opts.Level == nil {
		return true
	}
	return level >= h.opts.Level.Level()
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
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

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	fmt.Fprintf(h.out, "%s %s %s %s\n",
		timeStr,
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &PrettyHandler{
		out:   h.out,
		opts:  h.opts,
		attrs: newAttrs,
		group: h.group,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		out:   h.out,
		opts:  h.opts,
		attrs: h.attrs,
		group: name,
	}
}
