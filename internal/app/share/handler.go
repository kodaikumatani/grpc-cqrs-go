package share

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/identity"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/authz"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/grpcerr"
	pb "github.com/kodaikumatani/grpc-cqrs-go/pkg/pb/share"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	msgUnauthenticated = "unauthenticated"
	msgAlreadyShared   = "recipe already shared with this user"
	msgInvalidRelation = "invalid relation"
)

type handler struct {
	pb.UnimplementedShareServiceServer
	command *Command
}

func NewHandler(command *Command) pb.ShareServiceServer {
	return &handler{command: command}
}

func (h *handler) ShareRecipe(
	ctx context.Context,
	in *pb.ShareRecipeRequest,
) (*emptypb.Empty, error) {
	request := struct {
		RecipeId     string `validate:"required"`
		TargetUserId string `validate:"required"`
		Relation     string `validate:"required"`
	}{
		RecipeId:     in.GetRecipeId(),
		TargetUserId: in.GetTargetUserId(),
		Relation:     in.GetRelation(),
	}

	if err := validator.New().Struct(request); err != nil {
		return nil, grpcerr.WithStatus(err, codes.InvalidArgument, err.Error())
	}

	relation, err := authz.NewRelation(request.Relation)
	if err != nil {
		return nil, grpcerr.WithStatus(err, codes.InvalidArgument, msgInvalidRelation)
	}

	userID, err := identity.UserID(ctx)
	if err != nil {
		return nil, grpcerr.WithStatus(err, codes.Unauthenticated, msgUnauthenticated)
	}

	if err := h.command.ShareRecipe(ctx,
		userID.String(),
		request.RecipeId,
		request.TargetUserId,
		relation); err != nil {
		if errors.Is(err, authz.ErrAlreadyExists) {
			return nil, grpcerr.WithStatus(err, codes.AlreadyExists, msgAlreadyShared)
		}
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
