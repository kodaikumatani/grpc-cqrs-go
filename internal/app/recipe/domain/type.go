package domain

import (
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
}

// NewRecipe は業務フィールドを受け取ってレシピを構築する。
// timestamps は永続化層(DB の default / now())が管理するため domain は持たない。
func NewRecipe(
	id uuid.UUID,
	userID ulid.ULID,
	title, description string,
	visibility Visibility,
) *Recipe {
	return &Recipe{
		id:          id,
		userID:      userID,
		title:       title,
		description: description,
		visibility:  visibility,
	}
}

func (r *Recipe) Update(title, description string) {
	r.title = title
	r.description = description
}

func (r *Recipe) ID() uuid.UUID          { return r.id }
func (r *Recipe) UserID() ulid.ULID      { return r.userID }
func (r *Recipe) Title() string          { return r.title }
func (r *Recipe) Description() string    { return r.description }
func (r *Recipe) Visibility() Visibility { return r.visibility }

type Visibility string

var (
	VisibilityPublic     Visibility = "public"
	VisibilityPrivate    Visibility = "private"
	VisibilityRestricted Visibility = "restricted"
)
