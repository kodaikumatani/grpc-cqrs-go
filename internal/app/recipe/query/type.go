package query

import (
	"time"

	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/entity"
)

type RecipeWithUser struct {
	ID          string
	UserID      string
	Title       string
	Description string
	Visibility  entity.Visibility
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserName    string
	UserEmail   string
}
