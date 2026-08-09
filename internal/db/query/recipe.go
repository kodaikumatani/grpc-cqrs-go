package query

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/domain"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/query"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/db/gen"
)

type recipe struct {
	queries *gen.Queries
}

func NewRecipe(pool *pgxpool.Pool) query.Storage {
	return &recipe{queries: gen.New(pool)}
}

func (r *recipe) Get(ctx context.Context, id uuid.UUID) (*query.RecipeWithUser, error) {
	row, err := r.queries.GetRecipeWithUser(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrRecipeNotFound
	}
	if err != nil {
		return nil, err
	}

	visibility, err := domain.NewVisibility(string(row.Visibility))
	if err != nil {
		return nil, err
	}

	return &query.RecipeWithUser{
		ID:          row.ID.String(),
		UserID:      row.UserID.String(),
		Title:       row.Title,
		Description: row.Description,
		Visibility:  visibility,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		UserName:    row.UserName,
		UserEmail:   row.UserEmail,
	}, nil
}
