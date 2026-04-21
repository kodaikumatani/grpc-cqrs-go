package share

import (
	"context"

	"github.com/google/uuid"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/authn"
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
	if err := u.checker.CanShareRecipe(ctx, recipeID); err != nil {
		return err
	}

	userID, ok := ctx.Value(authn.UIDKey{}).(ulid.ULID)
	if !ok {
		return authn.ErrUnauthenticated
	}

	tuple := authz.NewTuple(
		uuid.New(),
		authz.ObjectRecipe,
		recipeID,
		relation,
		userID,
	)

	//if err := u.storage.CreateTuple(ctx, *tuple) {
	//	return nil, err
	//}

	return &recipe, nil
}
