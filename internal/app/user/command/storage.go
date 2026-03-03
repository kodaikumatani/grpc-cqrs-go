package command

import (
	"context"

	"github.com/kodaikumatani/grpc-cqrs/internal/app/user/domain"
)

type Storage interface {
	Create(ctx context.Context, user *domain.User) error
}
