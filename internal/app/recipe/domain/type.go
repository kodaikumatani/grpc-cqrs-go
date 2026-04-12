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
	Visibility  Visibility
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Visibility string

var (
	VisibilityPublic     Visibility = "public"
	VisibilityPrivate    Visibility = "private"
	VisibilityRestricted Visibility = "restricted"
)

func (r *Recipe) Update(title, description string) {
	r.Title = title
	r.Description = description
	r.UpdatedAt = time.Now()
}
