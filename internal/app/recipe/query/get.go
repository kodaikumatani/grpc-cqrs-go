package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/domain"
)

type Query struct {
	storage Storage
}

func NewQuery(
	storage Storage,
) *Query {
	return &Query{
		storage: storage,
	}
}

func (q *Query) Get(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Recipe, error) {
	recipe, err := q.storage.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}
