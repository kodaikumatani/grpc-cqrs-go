package interceptor

import (
	"context"
	"strings"

	"github.com/kodaikumatani/grpc-cqrs-go/internal/grpcerr"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/identity"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/identity/gateway"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

const msgUnauthenticated = "unauthenticated"

// publicMethodPrefixes は認証不要のインフラ系メソッド。
// reflection は grpcurl 等の探索用、health は本番でも probe が無認証で叩ける必要がある。
var publicMethodPrefixes = []string{
	"/grpc.reflection.",
	"/grpc.health.",
}

// GatewayAuthUnaryInterceptor は前段(ESP 等)が付与した検証済みユーザー情報ヘッダから
// UID を取り出し ctx(identity.UIDKey)に格納する。トークン検証は前段が済ませている前提
// （app 内では検証しない）。cause は保持し、client には安全な Unauthenticated を返す。
func GatewayAuthUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		uid, err := userIDFromMetadata(ctx)
		if err != nil {
			return nil, grpcerr.WithStatus(err, codes.Unauthenticated, msgUnauthenticated)
		}

		return handler(context.WithValue(ctx, identity.UIDKey{}, uid), req)
	}
}

func isPublicMethod(fullMethod string) bool {
	for _, prefix := range publicMethodPrefixes {
		if strings.HasPrefix(fullMethod, prefix) {
			return true
		}
	}
	return false
}

func userIDFromMetadata(ctx context.Context) (ulid.ULID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ulid.ULID{}, identity.ErrUnauthenticated
	}

	vals := md.Get(gateway.HeaderUserInfo)
	if len(vals) == 0 {
		return ulid.ULID{}, identity.ErrUnauthenticated
	}

	return gateway.Parse(vals[0])
}
