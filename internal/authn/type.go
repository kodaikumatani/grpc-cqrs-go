package authn

import (
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
)

var ErrUnauthenticated = errors.New("user is not authenticated")

type Verifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (context.Context, error)
}

type UIDKey struct{}

func UserID(ctx context.Context) (ulid.ULID, bool) {
	id, ok := ctx.Value(UIDKey{}).(ulid.ULID)

	return id, ok
}
