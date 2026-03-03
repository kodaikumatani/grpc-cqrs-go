package command

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/command"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/domain"
	"github.com/kodaikumatani/grpc-cqrs/internal/db/gen"
)

type recipe struct {
	queries *gen.Queries
}

func NewRecipe(pool *pgxpool.Pool) command.Storage {
	return &recipe{queries: gen.New(pool)}
}

func (r *recipe) Create(ctx context.Context, rec *domain.Recipe) error {
	return r.queries.CreateRecipe(ctx, gen.CreateRecipeParams{
		ID:          rec.ID,
		UserID:      rec.UserID.String(),
		Title:       rec.Title,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	})
}
