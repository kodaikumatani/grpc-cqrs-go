package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/domain"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/query"
	"github.com/kodaikumatani/grpc-cqrs/internal/db/gen"
)

type recipe struct {
	queries *gen.Queries
}

func NewRecipe(pool *pgxpool.Pool) query.Storage {
	return &recipe{queries: gen.New(pool)}
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
