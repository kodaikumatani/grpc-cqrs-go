package db

import (
	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/db/recipe"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/db/tuple"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/db/user"
)

var Set = wire.NewSet(
	NewPool,
	recipe.Set,
	user.Set,
	tuple.Set,
)
