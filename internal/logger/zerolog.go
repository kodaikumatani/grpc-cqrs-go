package logger

import (
	"context"
	"io"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
)

func withLogger(ctx context.Context, out io.Writer) context.Context {
	logger := zerolog.
		New(out).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Stack().
		Logger()

	zerolog.DefaultContextLogger = &logger

	return logger.WithContext(ctx)
}

func WithLogger(ctx context.Context) context.Context {
	if fd := os.Stdout.Fd(); isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd) {
		return withLogger(ctx, zerolog.NewConsoleWriter())
	}

	return withLogger(ctx, io.Writer(os.Stdout))
}

func WithLevel(ctx context.Context, level string) (context.Context, error) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	return withLogger(ctx, zerolog.Ctx(ctx).Level(lvl)), nil
}
