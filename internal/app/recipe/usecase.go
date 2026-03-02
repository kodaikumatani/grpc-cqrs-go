package recipe

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type UseCase struct {
	storage Storage
}

func NewUseCase(
	storage Storage,
) *UseCase {
	return &UseCase{
		storage: storage,
	}
}

func (u *UseCase) Create(
	ctx context.Context,
	userID,
	title,
	description string,
) (*Recipe, error) {
	recipe := Recipe{
		ID:          lo.Must(uuid.NewV7()),
		UserID:      userID,
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
