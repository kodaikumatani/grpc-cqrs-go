package identity

import (
	"context"
	"errors"
)

var ErrUnauthenticated = errors.New("user is not authenticated")

type UIDKey struct{}

// UserID は認証済みユーザーの識別子を返す。値は IdP の subject(sub)で、
// ULID とは限らない不透明な文字列（例: Google OAuth は数値文字列）。
func UserID(ctx context.Context) (string, error) {
	id, ok := ctx.Value(UIDKey{}).(string)
	if !ok || id == "" {
		return "", ErrUnauthenticated
	}

	return id, nil
}
