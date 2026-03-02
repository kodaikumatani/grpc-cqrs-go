package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/kodaikumatani/grpc-cqrs/internal/interceptor"
	"github.com/kodaikumatani/grpc-cqrs/internal/logger"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func run(ctx context.Context) error {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}

	services, cleanup, err := initializeServices(ctx, os.Getenv("DATABASE_URL"))
	defer cleanup()
	if err != nil {
		return errors.Wrap(err, "failed to initialize services")
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryUnaryInterceptor(),
			interceptor.LoggingUnaryInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			interceptor.RecoveryStreamInterceptor(),
			interceptor.LoggingStreamInterceptor(),
		),
	)
	s = services.Register(s)
	reflection.Register(s)

	log.Ctx(ctx).Info().Msgf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to serve")
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx = logger.WithLogger(ctx)

	if err := run(ctx); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Send()
	}
}
