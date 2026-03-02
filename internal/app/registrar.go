package app

import (
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe"
	"google.golang.org/grpc"
)

type Registrar struct {
	recipeHandler *recipe.Handler
}

func NewRegistrar(
	recipeHandler *recipe.Handler,
) *Registrar {
	return &Registrar{
		recipeHandler: recipeHandler,
	}
}

func (r *Registrar) Register(app *grpc.Server) *grpc.Server {
	r.recipeHandler.RegisterService(app)

	return app
}
