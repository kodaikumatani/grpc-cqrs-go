package interceptor

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryUnaryInterceptor() grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(
		panicRecoveryHandler(),
	)
}

func RecoveryStreamInterceptor() grpc.StreamServerInterceptor {
	return recovery.StreamServerInterceptor(
		panicRecoveryHandler(),
	)
}

func panicRecoveryHandler() recovery.Option {
	return recovery.WithRecoveryHandlerContext(func(ctx context.Context, p any) error {
		log.Ctx(ctx).Error().
			Interface("panic", p).
			Msg("panic recovered in gRPC handler")

		return status.Error(codes.Internal, "internal server error")
	})
}
