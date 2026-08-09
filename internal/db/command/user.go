package command

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/user/command"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/user/domain"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/db/gen"
)

type user struct {
	queries *gen.Queries
}

func NewUser(pool *pgxpool.Pool) command.Storage {
	return &user{queries: gen.New(pool)}
}

func (u *user) Create(ctx context.Context, usr *domain.User) error {
	err := u.queries.CreateUser(ctx, gen.CreateUserParams{
		ID:    usr.ID(),
		Name:  usr.Name(),
		Email: usr.Email(),
	})

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrAlreadyExists
	}

	return err
}
