package main

import (
	"context"
	"fmt"
	"go-keeper/cmd/server/config"
	"go-keeper/internal/server"
	"go-keeper/internal/server/data/database"
	"go-keeper/internal/server/data/repository"
	"go-keeper/internal/server/services"
	"go-keeper/pkg/jwtfactory"
	"go-keeper/pkg/logging"
	"go-keeper/pkg/pgxstorage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logging.NewZapLogger(zapcore.DebugLevel)
	if err != nil {
		log.Fatal(err)
	}

	rootCtx, cancelCtx := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	defer cancelCtx()

	dbFactory := database.NewPgxDatabaseFactory(cfg.DB)
	storage, err := pgxstorage.New(dbFactory, cfg.DB.RetryAttemptDelays)
	if err != nil {
		log.Fatal(err)
	}
	authRepository := repository.NewAuthRepository(storage, logger)

	tokenFactory := jwtfactory.New(cfg.JWTConfig)

	authService := services.NewAuthService(authRepository, tokenFactory)

	srv := server.NewServer(cfg.ServerGRPC, authService)

	if err := run(rootCtx, cfg, srv, logger); err != nil {
		logger.ErrorCtx(rootCtx, "Server shutdown with error", zap.Error(err))
	} else {
		logger.InfoCtx(rootCtx, "Server shutdown gracefully")
	}
}

func run(
	rootCtx context.Context,
	cfg *config.Config,
	server *server.Server,
	logger *logging.ZapLogger,
) error {
	g, ctx := errgroup.WithContext(rootCtx)

	context.AfterFunc(ctx, func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancelCtx()

		<-ctx.Done()
		log.Fatal("failed to gracefully shutdown the server")
	})

	g.Go(func() error {
		if err := server.Run(); err != nil {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		defer logger.InfoCtx(ctx, "Shutting down server")
		<-ctx.Done()
		if err := server.Stop(); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("goroutine error occured: %w", err)
	}

	return nil
}
