package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/command"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/domain"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/db/gen"
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
		UserID:      rec.UserID,
		Title:       rec.Title,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	})
}

func (r *recipe) Get(ctx context.Context, id uuid.UUID) (*domain.Recipe, error) {
	row, err := r.queries.GetRecipe(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Recipe{
		ID:          row.ID,
		UserID:      row.UserID,
		Title:       row.Title,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func (r *recipe) Update(ctx context.Context, rec *domain.Recipe) error {
	return r.queries.UpdateRecipe(ctx, gen.UpdateRecipeParams{
		ID:          rec.ID,
		Title:       rec.Title,
		Description: rec.Description,
		UpdatedAt:   rec.UpdatedAt,
	})
}
