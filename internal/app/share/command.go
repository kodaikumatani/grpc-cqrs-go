package share

import (
	"context"

	"github.com/google/uuid"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/authz"
	"github.com/oklog/ulid/v2"
)

type Command struct {
	storage authz.Storage
	checker authz.Checker
}

func NewCommand(
	storage authz.Storage,
	checker authz.Checker,
) *Command {
	return &Command{
		storage: storage,
		checker: checker,
	}
}

func (u *Command) ShareRecipe(
	ctx context.Context,
	recipeID string,
	targetUserID string,
	relation string,
) error {
	uid, err := ulid.Parse(targetUserID)
	if err != nil {
		return err
	}

	if err := u.checker.CanShareRecipe(ctx, recipeID); err != nil {
		return err
	}

	rel, err := authz.NewRelation(relation)
	if err != nil {
		return err
	}

	tuple := authz.NewTuple(
		uuid.New(),
		authz.ObjectRecipe,
		recipeID,
		rel,
		uid,
	)

	if err := u.storage.CreateTuple(ctx, tuple); err != nil {
		return err
	}

	return nil
}
