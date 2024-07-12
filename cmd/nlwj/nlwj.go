package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"nljw/internal/api"
	"nljw/internal/api/spec"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/phenpessoa/gutils/netutils/httputils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("\nCiao ;)")
}

func run(ctx context.Context) error {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	logger = logger.Named("nlwj app")
	defer func() { _ = logger.Sync() }()

	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME")))
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	si := api.NewApi(pool, logger)

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Recoverer, httputils.ChiLogger(logger))
	r.Mount("/", spec.Handler(&si))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	defer func() {
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown server", zap.Error(err))
		}

	}()

	errChan := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <- errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}
