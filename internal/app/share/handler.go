package share

import (
	"context"

	pb "github.com/kodaikumatani/grpc-cqrs-go/pkg/pb/share"
)

type handler struct {
	pb.UnimplementedShareServiceServer
}

func NewHandler() pb.UnimplementedShareServiceServer {
	return &handler{}
}

func (h *handler) ShareRecipe(
	ctx context.Context,
	in *pb.ShareRecipeRequest,
) error {
	request := struct {
		RecipeId     string
		TargetUserId string
		Relation     string
	}{
		RecipeId:     in.GetRecipeId(),
		TargetUserId: in.GetTargetUserId(),
		Relation:     in.GetRelation(),
	}

	return nil
}
