package command

import (
	"context"

	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/domain"
)

type Storage interface {
	Create(ctx context.Context, recipe *domain.Recipe) error
}
