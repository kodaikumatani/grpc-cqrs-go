package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/domain"
)

type Storage interface {
	Create(ctx context.Context, recipe *domain.Recipe) error
	Get(ctx context.Context, id uuid.UUID) (*domain.Recipe, error)
	Update(ctx context.Context, rec *domain.Recipe) error
}
