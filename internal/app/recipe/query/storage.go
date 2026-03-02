package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/domain"
)

type Storage interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Recipe, error)
}
