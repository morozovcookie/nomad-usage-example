package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/morozovcookie/nomand-usage-example/cmd/grpcserver/config"
	"github.com/morozovcookie/nomand-usage-example/grpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/health"
	healthV1 "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	appname = "grpc-server"
)

func main() {
	logger, err := initLogger()
	if err != nil {
		log.Fatal("init zap-logger error", err)
	}

	cfg := config.New()
	if err = cfg.Parse(); err != nil {
		logger.Fatal("parse config error", zap.Error(err))
	}

	srv := grpc.NewServer(cfg.GRPCServerConfig.Address,
		logger.With(zap.String("appname", appname)))

	healthSrv := health.NewServer()
	healthV1.RegisterHealthServer(srv, healthSrv)

	logger.Info("starting application")

	eg, ctx := errgroup.WithContext(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		logger.Info("starting grpc server", zap.String("address", srv.Address()))

		if err := srv.ListenAndServe(); err != nil {
			return errors.Wrap(err, "grpc server error")
		}

		return nil
	})

	logger.Info("application started")

	select {
	case <-quit:
		break
	case <-ctx.Done():
		break
	}

	logger.Info("stopping application")

	healthSrv.Shutdown()

	srv.Close()

	if err := eg.Wait(); err != nil {
		logger.Error("wait error", zap.Error(err))
	}

	logger.Info("application stopped")
}

func initLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, errors.Wrap(err, "init production logger error")
	}

	logger = logger.With(zap.String("appname", appname))

	return logger, nil
}
