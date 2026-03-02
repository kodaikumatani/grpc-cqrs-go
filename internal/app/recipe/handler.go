package recipe

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/command"
	"github.com/kodaikumatani/grpc-cqrs/internal/app/recipe/query"
	pb "github.com/kodaikumatani/grpc-cqrs/pkg/pb/recipe"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type handler struct {
	pb.UnimplementedRecipeServiceServer
	command *command.Command
	query   *query.Query
}

func NewHandler(
	command *command.Command,
	query *query.Query,
) pb.RecipeServiceServer {
	return &handler{
		command: command,
		query:   query,
	}
}

func (h *handler) CreateRecipe(
	ctx context.Context,
	in *pb.CreateRecipeRequest,
) (*pb.CreateRecipeResponse, error) {
	request := struct {
		UserID      string `validate:"required"`
		Title       string `validate:"required"`
		Description string `validate:"required"`
	}{
		UserID:      in.GetUserId(),
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
	}

	if err := validator.New().Struct(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := h.command.Create(ctx,
		request.UserID,
		request.Title,
		request.Description)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateRecipeResponse{
		RecipeId: result.ID.String(),
	}, nil
}

func (h *handler) GetRecipe(
	ctx context.Context,
	in *pb.GetRecipeRequest,
) (*pb.GetRecipeResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := h.query.Get(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetRecipeResponse{
		Recipe: &pb.Recipe{
			Id:          result.ID.String(),
			UserId:      result.UserID,
			Title:       result.Title,
			Description: result.Description,
			CreatedAt:   timestamppb.New(result.CreatedAt),
			UpdatedAt:   timestamppb.New(result.UpdatedAt),
		},
	}, nil
}
