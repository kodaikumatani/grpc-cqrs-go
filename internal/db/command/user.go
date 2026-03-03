package command

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/user/command"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/user/domain"
	"github.com/kodaikumatani/grpc-cqrs/internal/db/gen"
)

type user struct {
	queries *gen.Queries
}

func NewUser(pool *pgxpool.Pool) command.Storage {
	return &user{queries: gen.New(pool)}
}

func (u *user) Create(ctx context.Context, usr *domain.User) error {
	return u.queries.CreateUser(ctx, gen.CreateUserParams{
		ID:        usr.ID,
		Name:      usr.Name,
		Email:     usr.Email,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
	})
}
