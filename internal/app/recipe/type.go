package recipe

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

var (
	ErrRecipeNotFound = errors.New("recipe not found")
)

type Recipe struct {
	ID          uuid.UUID
	UserID      string
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Storage interface {
	Create(ctx context.Context, recipe *Recipe) error
}
