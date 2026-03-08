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
