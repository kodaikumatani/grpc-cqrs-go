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
	return c.hasPermission(ctx, ObjectRecipe, recipeID, PermView)
}

func (c Checker) CanEditRecipe(ctx context.Context, recipeID string) error {
	return c.hasPermission(ctx, ObjectRecipe, recipeID, PermEdit)
}

func (c Checker) CanShareRecipe(ctx context.Context, recipeID string) error {
	return c.hasPermission(ctx, ObjectRecipe, recipeID, PermShare)
}

func (c Checker) IsRecipeOwner(ctx context.Context, recipeID string) error {
	return c.hasRelation(ctx, ObjectRecipe, recipeID, RelOwner)
}

func (c Checker) hasPermission(
	ctx context.Context,
	objectType ObjectType,
	objectID string,
	perm Permission,
) error {
	uid, ok := ctx.Value(authn.UIDKey{}).(string)
	if !ok {
		return authn.ErrUnauthenticated
	}

	userID, err := ulid.Parse(uid)
	if err != nil {
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

	for _, relation := range permissionRelations[objectType][perm] {
		for _, tuple := range tuples {
			if relation == tuple.Relation {
				return nil
			}
		}
	}

	return ErrPermissionDenied
}

func (c Checker) hasRelation(
	ctx context.Context,
	objectType ObjectType,
	objectID string,
	relation Relation,
) error {
	uid, ok := ctx.Value(authn.UIDKey{}).(string)
	if !ok {
		return authn.ErrUnauthenticated
	}

	userID, err := ulid.Parse(uid)
	if err != nil {
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

	for _, tuple := range tuples {
		if relation == tuple.Relation {
			return nil
		}
	}

	return ErrPermissionDenied
}
