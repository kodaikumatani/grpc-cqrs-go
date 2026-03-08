package db

import (
	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/db/command"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/db/query"
)

var Set = wire.NewSet(
	NewPool,
	command.Set,
	query.Set,
)
