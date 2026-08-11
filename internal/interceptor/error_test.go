package interceptor

import (
	"testing"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
)

func TestLevelForCode(t *testing.T) {
	tests := []struct {
		code codes.Code
		want zerolog.Level
	}{
		// 想定内のクライアントエラーは info(warn にしない)
		{codes.PermissionDenied, zerolog.InfoLevel},
		{codes.InvalidArgument, zerolog.InfoLevel},
		{codes.AlreadyExists, zerolog.InfoLevel},
		{codes.Unauthenticated, zerolog.InfoLevel},
		// サーバ起因は error
		{codes.Internal, zerolog.ErrorLevel},
		{codes.Unknown, zerolog.ErrorLevel},
		// NotFound と一時的な障害は warn
		{codes.NotFound, zerolog.WarnLevel},
		{codes.Unavailable, zerolog.WarnLevel},
		{codes.DeadlineExceeded, zerolog.WarnLevel},
	}

	for _, tt := range tests {
		if got := levelForCode(tt.code); got != tt.want {
			t.Errorf("levelForCode(%v) = %v, want %v", tt.code, got, tt.want)
		}
	}
}
