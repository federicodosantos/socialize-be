package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/federicodosantos/socialize/internal/app"
	"github.com/federicodosantos/socialize/pkg/database/mysql"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type operation func(context.Context) error

func main() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()
	sugar := logger.Sugar()

	err := godotenv.Load()
	if err != nil {
		sugar.Fatalf("cannot load env file: %v", err)
	}

	db := mysql.DBInit()
	defer db.Close()

	router := chi.NewRouter()

	bootstrap := app.NewBootstrap(db, router, sugar)
	bootstrap.InitApp()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8060"),
		Handler: router,
	}

	go func() {
		sugar.Infof("server is running on port %s", "8060")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("server failed to start: %v", err)
		}
	}()

	wait := gracefullyShutdown(context.Background(), 5*time.Second, sugar, map[string]operation{
		"http-server": func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
		"db-user": func(ctx context.Context) error {
			return db.Close()
		},
	})

	<-wait
	sugar.Info("application stopped gracefully")
}

func gracefullyShutdown(ctx context.Context, timeout time.Duration, sugar *zap.SugaredLogger, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		// Blok sampai menerima sinyal
		<-c
		sugar.Info("shutting down application")

		// Timeout untuk proses shutdown
		timeoutFunc := time.AfterFunc(timeout, func() {
			sugar.Warnf("timeout reached (%v), forcing exit", timeout)
			os.Exit(1)
		})
		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Proses shutdown asinkron
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()
				sugar.Infof("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					sugar.Errorf("%s cleanup failed: %v", innerKey, err)
					return
				}
				sugar.Infof("%s shut down gracefully", innerKey)
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}
