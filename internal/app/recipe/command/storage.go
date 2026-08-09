package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/entity"
)

type Storage interface {
	Create(ctx context.Context, recipe *entity.Recipe) error
	Get(ctx context.Context, id uuid.UUID) (*entity.Recipe, error)
	Update(ctx context.Context, rec *entity.Recipe) error
}
