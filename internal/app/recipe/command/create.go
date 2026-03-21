package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/domain"
	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
)

type Command struct {
	storage Storage
}

func NewCommand(
	storage Storage,
) *Command {
	return &Command{
		storage: storage,
	}
}

func (u *Command) Create(
	ctx context.Context,
	userID,
	title,
	description string,
) (*domain.Recipe, error) {
	uid, err := ulid.Parse(userID)
	if err != nil {
		return nil, err
	}

	recipe := domain.Recipe{
		ID:          lo.Must(uuid.NewV7()),
		UserID:      uid,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.storage.Create(ctx, &recipe); err != nil {
		return nil, err
	}

	return &recipe, nil
}

func (u *Command) Update(
	ctx context.Context,
	id uuid.UUID,
	title, description string,
) error {
	recipe, err := u.storage.Get(ctx, id)
	if err != nil {
		return err
	}

	recipe.Update(title, description)

	if err := u.storage.Update(ctx, recipe); err != nil {
		return err
	}

	return nil
}
