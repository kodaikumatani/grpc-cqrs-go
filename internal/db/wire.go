package db

import (
	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs/internal/db/recipe"
)

var Set = wire.NewSet(
	NewPool,
	recipe.NewRepository,
)
