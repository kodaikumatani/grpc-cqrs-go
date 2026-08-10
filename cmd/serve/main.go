package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/interceptor"
	"github.com/kodaikumatani/grpc-cqrs-go/internal/logger"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
	// host はデフォルトで loopback。app はヘッダ(x-endpoint-api-userinfo)を信頼するため、
	// 前段(同一 Pod の Envoy 等)からのみ到達させる。別コンテナから叩く構成では 0.0.0.0 に上書き。
	host = flag.String("host", "127.0.0.1", "The server listen address")
)

func run(ctx context.Context) error {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
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
			interceptor.LoggingUnaryInterceptor(),
			interceptor.RecoveryUnaryInterceptor(),
			interceptor.ErrorHandlingUnaryInterceptor(),
			interceptor.GatewayAuthUnaryInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			interceptor.LoggingStreamInterceptor(),
			interceptor.RecoveryStreamInterceptor(),
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
