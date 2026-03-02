package app

import (
	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe"
)

var Set = wire.NewSet(
	recipe.NewUseCase,
	recipe.NewHandler,
	NewRegistrar,
)
