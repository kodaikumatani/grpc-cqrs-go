package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/domain"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/authz"
)

type Query struct {
	storage Storage
	checker authz.Checker
}

func NewQuery(
	storage Storage,
	checker authz.Checker,
) *Query {
	return &Query{
		storage: storage,
		checker: checker,
	}
}

func (q *Query) Get(
	ctx context.Context,
	id uuid.UUID,
) (*RecipeWithUser, error) {
	result, err := q.storage.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	switch result.Visibility {
	case domain.VisibilityPublic:
		// 誰でも可（素通り）
	case domain.VisibilityPrivate:
		if err := q.checker.IsRecipeOwner(ctx, id.String()); err != nil {
			return nil, err
		}
	case domain.VisibilityRestricted:
		if err := q.checker.CanViewRecipe(ctx, id.String()); err != nil {
			return nil, err
		}
	default:
		return nil, authz.ErrPermissionDenied
	}

	return result, nil
}
