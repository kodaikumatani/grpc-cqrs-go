package share

import (
	"context"

	"github.com/go-playground/validator/v10"
	pb "github.com/kodaikumatani/grpc-cqrs-go/pkg/pb/share"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.command.ShareRecipe(ctx,
		request.RecipeId,
		request.TargetUserId,
		request.Relation); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
