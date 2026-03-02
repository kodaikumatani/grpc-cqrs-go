package db

import (
	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs/internal/db/command"
	"github.com/kodaikumatani/grpc-cqrs/internal/db/query"
)

var Set = wire.NewSet(
	NewPool,
	command.Set,
	query.Set,
)
