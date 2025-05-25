package main

import (
	"context"
	"fmt"
	"go-keeper/cmd/client/config"
	"go-keeper/internal/client"
	"go-keeper/internal/client/data/repositories"
	"go-keeper/internal/client/data/storage"
	"go-keeper/internal/client/logic/commands"
	"go-keeper/internal/client/logic/requester"
	"go-keeper/internal/client/logic/requester/options"
	"go-keeper/internal/client/logic/services"
	"go-keeper/pkg/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

	logger, err := logging.NewZapLogger(zapcore.DebugLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	rootCtx, cancelCtx := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	defer cancelCtx()

	str, err := storage.NewFileStorage(cfg.LocalStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	tokenRepository := repositories.NewTokenRepository(str)
	cmds := commands.NewCommands(os.Stdin, os.Stdout)

	authReq := requester.New("localhost:8080", []options.Option{})
	authService := services.NewAuthService(authReq)

	storageReq := requester.New(
		"localhost:8080",
		[]options.Option{
			options.NewAuthOption(tokenRepository),
		},
	)
	storageService := services.NewStorageService(storageReq)

	cli := client.New(
		cfg.Client,
		tokenRepository,
		cmds,
		authService,
		storageService,
	)

	if err := run(rootCtx, cfg, cli, logger); err != nil {
		logger.ErrorCtx(rootCtx, "Client shutdown with error", zap.Error(err))
	}
}

func run(
	rootCtx context.Context,
	cfg *config.Config,
	cli *client.Client,
	logger *logging.ZapLogger,
) error {
	g, ctx := errgroup.WithContext(rootCtx)
	cctx, cancel := context.WithCancel(ctx)
	ctx = cctx

	context.AfterFunc(ctx, func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancelCtx()

		<-ctx.Done()
		log.Fatal("failed to gracefully shutdown the client")
	})

	g.Go(func() error {
		defer cancel()
		if err := cli.Run(ctx); err != nil {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		defer logger.InfoCtx(ctx, "Shutting down client")
		<-ctx.Done()
		if err := cli.Stop(); err != nil {
			return fmt.Errorf("failed to shutdown client: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("goroutine error occured: %w", err)
	}

	return nil
}
