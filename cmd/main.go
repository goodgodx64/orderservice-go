package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/goodgodx64/orderservice-go/internal/service"
	pb "github.com/goodgodx64/orderservice-go/pkg/api/grpc"
	"github.com/goodgodx64/orderservice-go/pkg/logger"

	"github.com/goodgodx64/orderservice-go/internal/config"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	ctx, err := logger.New(ctx)
	if err != nil {
		log.Fatal("failed to create logger", err)
	}

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to load config", zap.Error(err))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v\n", lis.Addr())
	srv := service.New()
	server := grpc.NewServer(grpc.UnaryInterceptor(logger.Interceptor))
	pb.RegisterOrderServiceServer(server, srv)

	if err := server.Serve(lis); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to serve", zap.Error(err))
	}
}
