package command

import (
	"context"
	"time"

	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/user/domain"
	"github.com/oklog/ulid/v2"
)

type Command struct {
	storage Storage
}

func NewCommand(storage Storage) *Command {
	return &Command{storage: storage}
}

func (c *Command) Create(
	ctx context.Context,
	name,
	email string,
) (*domain.User, error) {
	user := domain.User{
		ID:        ulid.Make(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := c.storage.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
