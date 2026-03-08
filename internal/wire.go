package internal

import (
	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/db"
)

var Set = wire.NewSet(
	app.Set,
	db.Set,
)
