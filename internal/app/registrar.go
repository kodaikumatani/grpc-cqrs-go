package app

import (
	recipepb "github.com/kodaikumatani/grpc-cqrs-go/pkg/pb/recipe"
	sharepb "github.com/kodaikumatani/grpc-cqrs-go/pkg/pb/share"
	userpb "github.com/kodaikumatani/grpc-cqrs-go/pkg/pb/user"
	"google.golang.org/grpc"
)

type Registrar struct {
	recipeHandler recipepb.RecipeServiceServer
	userHandler   userpb.UserServiceServer
	shareHandler  sharepb.ShareServiceServer
}

func NewRegistrar(
	recipeHandler recipepb.RecipeServiceServer,
	userHandler userpb.UserServiceServer,
	shareHandler sharepb.ShareServiceServer,
) *Registrar {
	return &Registrar{
		recipeHandler: recipeHandler,
		userHandler:   userHandler,
		shareHandler:  shareHandler,
	}
}

func (r *Registrar) Register(app *grpc.Server) *grpc.Server {
	recipepb.RegisterRecipeServiceServer(app, r.recipeHandler)
	userpb.RegisterUserServiceServer(app, r.userHandler)
	sharepb.RegisterShareServiceServer(app, r.shareHandler)

	return app
}
