package user

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/user/command"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/authn"
	pb "github.com/kodaikumatani/grpc-cqrs-go/pkg/pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	pb.UnimplementedUserServiceServer
	command *command.Command
}

func NewHandler(command *command.Command) pb.UserServiceServer {
	return &handler{command: command}
}

func (h *handler) CreateUser(
	ctx context.Context,
	in *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {
	request := struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}{
		Name:  in.GetName(),
		Email: in.GetEmail(),
	}

	if err := validator.New().Struct(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, ok := authn.UserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, authn.ErrUnauthenticated.Error())
	}

	result, err := h.command.Create(ctx, id, request.Name, request.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateUserResponse{
		UserId: result.ID().String(),
	}, nil
}
