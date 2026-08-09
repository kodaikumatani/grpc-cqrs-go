package query

import (
	"context"

	"github.com/google/uuid"
)

type Storage interface {
	Get(ctx context.Context, id uuid.UUID) (*RecipeWithUser, error)
}
