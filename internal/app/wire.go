package app

import (
	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/command"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/query"
)

var Set = wire.NewSet(
	command.NewCommand,
	query.NewQuery,
	recipe.NewHandler,
	NewRegistrar,
)
