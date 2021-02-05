package main

import (
	"context"
	"log"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/morozovcookie/nomand-usage-example/cmd/httpserver/config"
	"github.com/morozovcookie/nomand-usage-example/http"
	v1 "github.com/morozovcookie/nomand-usage-example/http/v1"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	appname = "http-server"

	shutdownTimeout = time.Second * 5

	rootPathPrefix = "/"
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

	var opts []http.ServerOption
	opts = initServerHandlers(opts, logger)

	srv := http.NewServer(cfg.HTTPServerConfig.Address,
		logger.With(zap.String("component", "http_server")),
		opts...)

	logger.Info("starting application")

	eg, ctx := errgroup.WithContext(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		logger.Info("starting http server", zap.String("address", srv.Address()))

		if err := srv.ListenAndServe(); err != nil {
			return errors.Wrap(err, "http server error")
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

	shutdownCtx, shutdownCancel := context.WithDeadline(context.Background(), time.Now().Add(shutdownTimeout))
	defer shutdownCancel()

	if err := srv.Close(shutdownCtx); err != nil {
		logger.Error("shutdown http server error", zap.Error(err))
	}

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

func initHealthzHandler(logger *zap.Logger) (string, stdhttp.Handler) {
	var (
		handlerConfig = &HealthHandlerConfiguration{
			LoggerInstance: logger.With(zap.String("component", "healthz_handler")),
		}
		handler = v1.NewHealthzHandler(handlerConfig)
	)

	return rootPathPrefix, handler
}

func initServerHandlers(opts []http.ServerOption, logger *zap.Logger) []http.ServerOption {
	{
		prefix, handler := initHealthzHandler(logger)
		opts = append(opts, http.WithHandler(prefix, handler))
	}

	return opts
}
