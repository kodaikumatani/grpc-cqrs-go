package command

import (
	"context"

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
	id ulid.ULID,
	name,
	email string,
) (*domain.User, error) {
	user := domain.NewUser(id, name, email)

	if err := c.storage.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
