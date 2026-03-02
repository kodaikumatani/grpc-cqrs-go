package recipe

import (
	"context"

	"github.com/go-playground/validator/v10"
	pb "github.com/kodaikumatani/grpc-cqrs/pkg/pb/recipe"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedRecipeServiceServer
	usecase *UseCase
}

func NewHandler(
	create *UseCase,
) *Handler {
	return &Handler{
		usecase: create,
	}
}

func (r *Handler) RegisterService(app *grpc.Server) {
	pb.RegisterRecipeServiceServer(app, r)
}

type createRecipeRequest struct {
	UserID      string `validate:"required"`
	Title       string `validate:"required"`
	Description string `validate:"required"`
}

func (h *Handler) CreateRecipe(
	ctx context.Context,
	in *pb.CreateRecipeRequest,
) (*pb.CreateRecipeResponse, error) {
	request := &createRecipeRequest{
		UserID:      in.GetUserId(),
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
	}

	if err := validator.New().Struct(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := h.usecase.Create(ctx,
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
