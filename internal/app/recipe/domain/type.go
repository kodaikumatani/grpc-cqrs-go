package domain

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

var (
	ErrRecipeNotFound = errors.New("recipe not found")
)

type Recipe struct {
	ID          uuid.UUID
	UserID      ulid.ULID
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *Recipe) Update(title, description string) {
	r.Title = title
	r.Description = description
	r.UpdatedAt = time.Now()
}
