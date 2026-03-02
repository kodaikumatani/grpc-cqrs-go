package internal

import (
	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs/internal/app"
	"github.com/kodaikumatani/grpc-cqrs/internal/db"
)

var Set = wire.NewSet(
	app.Set,
	db.Set,
)
