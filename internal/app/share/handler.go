package share

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/authn"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/grpcerr"
	pb "github.com/kodaikumatani/grpc-cqrs-go/pkg/pb/share"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	msgUnauthenticated   = "unauthenticated"
	msgUserAlreadyExists = "user already exists"
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

	userID, err := authn.UserID(ctx)
	if err != nil {
		return nil, grpcerr.WithStatus(err, codes.Unauthenticated, msgUnauthenticated)
	}

	if err := h.command.ShareRecipe(ctx,
		userID.String(),
		request.RecipeId,
		request.TargetUserId,
		request.Relation); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
