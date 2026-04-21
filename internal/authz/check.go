package authz

import (
	"context"

	"github.com/kodaikumatani/grpc-cqrs-go/internal/authn"
	"github.com/oklog/ulid/v2"
)

type Checker struct {
	storage Storage
}

func NewChecker(storage Storage) Checker {
	return Checker{storage: storage}
}

func (c Checker) CanViewRecipe(ctx context.Context, recipeID string) error {
	return c.check(ctx, ObjectRecipe, recipeID, PermViewRecipe)
}

func (c Checker) CanEditRecipe(ctx context.Context, recipeID string) error {
	return c.check(ctx, ObjectRecipe, recipeID, PermEditRecipe)
}

func (c Checker) CanShareRecipe(ctx context.Context, recipeID string) error {
	return c.check(ctx, ObjectRecipe, recipeID, PermShareRecipe)
}

func (c Checker) check(
	ctx context.Context,
	objectType ObjectType,
	objectID string,
	perm Permission,
) error {
	userID, ok := ctx.Value(authn.UIDKey{}).(ulid.ULID)
	if !ok {
		return authn.ErrUnauthenticated
	}

	tuples, err := c.storage.ListRelations(
		ctx,
		objectType,
		objectID,
		userID,
	)
	if err != nil {
		return err
	}

	for _, relation := range perm {
		for _, tuple := range tuples {
			if relation == tuple.Relation {
				return nil
			}
		}
	}

	return ErrPermissionDenied
}
