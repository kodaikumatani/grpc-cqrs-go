package command

import (
	"context"

	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/user/entity"
)

type Storage interface {
	Create(ctx context.Context, user *entity.User) error
}
