package interceptor

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func LoggingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(
		loggingAdapter(),
	)
}

func LoggingStreamInterceptor() grpc.StreamServerInterceptor {
	return logging.StreamServerInterceptor(
		loggingAdapter(),
	)
}

func loggingAdapter() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		logger := log.Ctx(ctx).With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			logger.Debug().Msg(msg)
		case logging.LevelInfo:
			logger.Info().Msg(msg)
		case logging.LevelWarn:
			logger.Warn().Msg(msg)
		case logging.LevelError:
			logger.Error().Msg(msg)
		default:
			logger.Info().Msg(msg)
		}
	})
}
