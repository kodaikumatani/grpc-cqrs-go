//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/kodaikumatani/grpc-cqrs-go/internal"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app"
)

type services struct {
	*app.Registrar
}

var set = wire.NewSet(
	internal.Set,
	wire.Struct(new(services), "*"),
)

func initializeServices(ctx context.Context, dsn string) (*services, func(), error) {
	panic(wire.Build(set))
}
