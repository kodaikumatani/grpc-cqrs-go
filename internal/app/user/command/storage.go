package command

import (
	"context"

	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/user/domain"
)

type Storage interface {
	Create(ctx context.Context, user *domain.User) error
}
