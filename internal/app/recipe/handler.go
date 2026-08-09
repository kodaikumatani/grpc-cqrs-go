package recipe

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/command"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/domain"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/recipe/query"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/authn"
	pb "github.com/kodaikumatani/grpc-cqrs-go/pkg/pb/recipe"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// toDomainVisibility は pb の enum を domain の Visibility に変換する。
// UNSPECIFIED / 未知の値は ErrInvalidVisibility を返す。
func toDomainVisibility(v pb.Visibility) (domain.Visibility, error) {
	switch v {
	case pb.Visibility_VISIBILITY_PUBLIC:
		return domain.VisibilityPublic, nil
	case pb.Visibility_VISIBILITY_PRIVATE:
		return domain.VisibilityPrivate, nil
	case pb.Visibility_VISIBILITY_RESTRICTED:
		return domain.VisibilityRestricted, nil
	default:
		return domain.Visibility{}, domain.ErrInvalidVisibility
	}
}

func toPBVisibility(v domain.Visibility) pb.Visibility {
	switch v {
	case domain.VisibilityPublic:
		return pb.Visibility_VISIBILITY_PUBLIC
	case domain.VisibilityPrivate:
		return pb.Visibility_VISIBILITY_PRIVATE
	case domain.VisibilityRestricted:
		return pb.Visibility_VISIBILITY_RESTRICTED
	default:
		return pb.Visibility_VISIBILITY_UNSPECIFIED
	}
}

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
		Title       string `validate:"required"`
		Description string `validate:"required"`
	}{
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
	}

	if err := validator.New().Struct(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, ok := ctx.Value(authn.UIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, ErrUserNotFound.Error())
	}

	result, err := h.command.Create(ctx,
		userID,
		request.Title,
		request.Description)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateRecipeResponse{
		RecipeId: result.ID().String(),
	}, nil
}

func (h *handler) UpdateRecipe(
	ctx context.Context,
	in *pb.UpdateRecipeRequest,
) (*pb.UpdateRecipeResponse, error) {
	request := struct {
		ID          string `validate:"required,uuid"`
		Title       string `validate:"required"`
		Description string `validate:"required"`
	}{
		ID:          in.GetId(),
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
	}

	if err := validator.New().Struct(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	recipeID := uuid.MustParse(request.ID)

	if err := h.command.Update(ctx, recipeID, request.Title, request.Description); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateRecipeResponse{Success: true}, nil
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
			Id:          result.ID,
			UserId:      result.UserID,
			Title:       result.Title,
			Description: result.Description,
			Visibility:  toPBVisibility(result.Visibility),
			CreatedAt:   timestamppb.New(result.CreatedAt),
			UpdatedAt:   timestamppb.New(result.UpdatedAt),
		},
		User: &pb.User{
			Id:    result.UserID,
			Name:  result.UserName,
			Email: result.UserEmail,
		},
	}, nil
}

func (h *handler) ChangeVisibility(
	ctx context.Context,
	in *pb.ChangeVisibilityRequest,
) (*pb.ChangeVisibilityResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	visibility, err := toDomainVisibility(in.GetVisibility())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.command.UpdateVisibility(ctx, id, visibility); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ChangeVisibilityResponse{Success: true}, nil
}
