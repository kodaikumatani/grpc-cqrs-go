package authz

import (
	"context"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type Storage interface {
	CreateTuple(ctx context.Context, tuple Tuple) error
	DeleteTuple(ctx context.Context, id uuid.UUID) error
	ListRelations(
		ctx context.Context,
		objectType ObjectType,
		objectID string,
		userID ulid.ULID,
	) ([]*Tuple, error)
}
