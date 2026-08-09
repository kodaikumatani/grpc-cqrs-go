package authz

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/authz"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/db/gen"
)

type tuple struct {
	queries *gen.Queries
}

func NewTuple(pool *pgxpool.Pool) authz.Storage {
	return &tuple{queries: gen.New(pool)}
}

func (t *tuple) CreateTuple(
	ctx context.Context,
	tuple *authz.Tuple,
) error {
	return t.queries.CreateTuple(ctx, gen.CreateTupleParams{
		ID:         uuid.New(),
		ObjectType: tuple.ObjectType().String(),
		ObjectID:   tuple.ObjectID(),
		Relation:   tuple.Relation().String(),
		UserID:     tuple.UserID(),
	})
}

func (t *tuple) DeleteTuple(
	ctx context.Context,
	id uuid.UUID,
) error {
	return t.queries.DeleteTuple(ctx, id)
}

func (t *tuple) ListRelations(
	ctx context.Context,
	objectType authz.ObjectType,
	objectID string,
	userID string,
) ([]*authz.Tuple, error) {
	rows, err := t.queries.ListRelations(ctx, gen.ListRelationsParams{
		ObjectType: objectType.String(),
		ObjectID:   objectID,
		UserID:     userID,
	})
	if err != nil {
		return nil, err
	}

	tuples := make([]*authz.Tuple, len(rows))
	for i, row := range rows {
		objectType, err := authz.NewObjectType(row.ObjectType)
		if err != nil {
			return nil, err
		}

		relation, err := authz.NewRelation(row.Relation)
		if err != nil {
			return nil, err
		}

		tuples[i] = authz.NewTuple(
			row.ID,
			objectType,
			row.ObjectID,
			relation,
			row.UserID,
		)
	}

	return tuples, nil
}
