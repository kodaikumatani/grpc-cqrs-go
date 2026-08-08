package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/domain"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/authz"
	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
)

type Command struct {
	storage Storage
	tuple   authz.Storage
	checker authz.Checker
}

func NewCommand(
	storage Storage,
	tuple authz.Storage,
	checker authz.Checker,
) *Command {
	return &Command{
		storage: storage,
		tuple:   tuple,
		checker: checker,
	}
}

func (u *Command) Create(
	ctx context.Context,
	userID,
	title,
	description string,
) (*domain.Recipe, error) {
	uid, err := ulid.Parse(userID)
	if err != nil {
		return nil, err
	}

	recipe := domain.NewRecipe(
		lo.Must(uuid.NewV7()),
		uid,
		title,
		description,
		domain.VisibilityPrivate,
	)

	tp := authz.NewTuple(
		uuid.New(),
		authz.ObjectRecipe,
		recipe.ID().String(),
		authz.RelOwner,
		uid,
	)

	if err := u.storage.Create(ctx, recipe); err != nil {
		return nil, err
	}

	if err := u.tuple.CreateTuple(ctx, tp); err != nil {
		return nil, err
	}

	return recipe, nil
}

func (u *Command) Update(
	ctx context.Context,
	id uuid.UUID,
	title, description string,
) error {
	if err := u.checker.
		CanEditRecipe(ctx, id.String()); err != nil {
		return err
	}

	recipe, err := u.storage.Get(ctx, id)
	if err != nil {
		return err
	}

	recipe.Update(title, description)

	if err := u.storage.Update(ctx, recipe); err != nil {
		return err
	}

	return nil
}
