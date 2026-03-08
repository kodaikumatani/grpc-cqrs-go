package user

import (
	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/user/command"
)

var Set = wire.NewSet(
	NewHandler,
	command.NewCommand,
)
