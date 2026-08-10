package internal

import (
	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/authz"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/db"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/health"
)

var Set = wire.NewSet(
	app.Set,
	db.Set,
	authz.Set,
	health.Set,
)
