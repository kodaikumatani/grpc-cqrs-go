package app

import (
	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/share"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/user"
)

var Set = wire.NewSet(
	NewRegistrar,
	recipe.Set,
	share.Set,
	user.Set,
)
