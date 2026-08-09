package interceptor

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const msgInternalServerError = "internal server error"

// ErrorHandlingUnaryInterceptor は handler が返したエラーを一元的に処理する。
//   - handler が付けた status（既知エラー）: cause を warn ログに残し、安全な status を返す
//   - 未変換の生エラー（想定外の内部エラー）: cause を error ログに残し、client には汎用文言
//
// これにより内部エラーの詳細を client に漏らさず、ログには常に cause を残せる。
func ErrorHandlingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		if st, ok := status.FromError(err); ok {
			log.Ctx(ctx).Warn().Err(err).Str("method", info.FullMethod).Msg("request failed")
			return resp, st.Err()
		}

		log.Ctx(ctx).Error().Err(err).Str("method", info.FullMethod).Msg("unhandled error")
		return resp, status.Error(codes.Internal, msgInternalServerError)
	}
}
