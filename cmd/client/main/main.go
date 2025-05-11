package main

import (
	"context"
	"fmt"
	"go-keeper/cmd/client/config"
	"go-keeper/internal/client"
	"go-keeper/pkg/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

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

	cli := client.New(cfg.Client)

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
		log.Fatal("failed to gracefully shutdown the server")
	})

	g.Go(func() error {
		defer cancel()
		if err := cli.Run(); err != nil {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		defer logger.InfoCtx(ctx, "Shutting down server")
		<-ctx.Done()
		if err := cli.Stop(); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("goroutine error occured: %w", err)
	}

	return nil
}
