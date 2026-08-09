package user

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/user/command"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/app/user/domain"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/authn"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/grpcerr"
	pb "github.com/kodaikumatani/grpc-cqrs-go/pkg/pb/user"
	"google.golang.org/grpc/codes"
)

const (
	msgUnauthenticated   = "unauthenticated"
	msgUserAlreadyExists = "user already exists"
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
		return nil, grpcerr.WithStatus(err, codes.InvalidArgument, err.Error())
	}

	id, err := authn.UserID(ctx)
	if err != nil {
		return nil, grpcerr.WithStatus(err, codes.Unauthenticated, msgUnauthenticated)
	}

	result, err := h.command.Create(ctx, id, request.Name, request.Email)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, grpcerr.WithStatus(err, codes.AlreadyExists, msgUserAlreadyExists)
		}
		return nil, err
	}

	return &pb.CreateUserResponse{
		UserId: result.ID().String(),
	}, nil
}
