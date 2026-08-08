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
	id          uuid.UUID
	userID      ulid.ULID
	title       string
	description string
	visibility  Visibility
	createdAt   time.Time
	updatedAt   time.Time
}

// NewRecipe は全フィールドを受け取ってレシピを構築する。
// id / visibility / timestamps は呼び出し側（app 層）が用意する。
func NewRecipe(
	id uuid.UUID,
	userID ulid.ULID,
	title, description string,
	visibility Visibility,
	createdAt, updatedAt time.Time,
) *Recipe {
	return &Recipe{
		id:          id,
		userID:      userID,
		title:       title,
		description: description,
		visibility:  visibility,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (r *Recipe) Update(title, description string) {
	r.title = title
	r.description = description
	r.updatedAt = time.Now()
}

func (r *Recipe) ID() uuid.UUID          { return r.id }
func (r *Recipe) UserID() ulid.ULID      { return r.userID }
func (r *Recipe) Title() string          { return r.title }
func (r *Recipe) Description() string    { return r.description }
func (r *Recipe) Visibility() Visibility { return r.visibility }
func (r *Recipe) CreatedAt() time.Time   { return r.createdAt }
func (r *Recipe) UpdatedAt() time.Time   { return r.updatedAt }

type Visibility string

var (
	VisibilityPublic     Visibility = "public"
	VisibilityPrivate    Visibility = "private"
	VisibilityRestricted Visibility = "restricted"
)
