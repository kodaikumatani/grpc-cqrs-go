package recipe

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe"
	"github.com/kodaikumatani/grpc-cqrs/internal/db/gen"
)

type repository struct {
	queries *gen.Queries
}

func NewRepository(pool *pgxpool.Pool) recipe.Storage {
	return &repository{queries: gen.New(pool)}
}

func (r *repository) Create(ctx context.Context, rec *recipe.Recipe) error {
	return r.queries.CreateRecipe(ctx, gen.CreateRecipeParams{
		ID:          rec.ID,
		UserID:      rec.UserID,
		Title:       rec.Title,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	})
}
