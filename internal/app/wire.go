package app

import (
	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe"
	recipecommand "github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/command"
	recipequery "github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/query"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/user"
	usercommand "github.com/kodaikumatani/grpc-cqrs/internal/app/user/command"
)

var Set = wire.NewSet(
	recipecommand.NewCommand,
	recipequery.NewQuery,
	recipe.NewHandler,
	usercommand.NewCommand,
	user.NewHandler,
	NewRegistrar,
)
