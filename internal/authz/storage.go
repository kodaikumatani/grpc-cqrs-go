package authz

import (
	"context"

	"github.com/google/uuid"
)

type Storage interface {
	CreateTuple(ctx context.Context, tuple *Tuple) error
	DeleteTuple(ctx context.Context, id uuid.UUID) error
	ListRelations(
		ctx context.Context,
		objectType ObjectType,
		objectID string,
		userID string,
	) ([]*Tuple, error)
}
