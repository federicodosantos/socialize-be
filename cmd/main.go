package main

import (
	"context"
	"fmt"
	"log"
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
)

type operation func(context.Context) error

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("cannot load env file due to %s", err)
	}

	db := mysql.DBInit()
	defer db.Close()

	router := chi.NewRouter()

	bootstrap := app.NewBootstrap(db, router)

	bootstrap.InitApp()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler: router,
	}

	go func() {
		log.Printf("server is running on port %s", os.Getenv("APP_PORT"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("cannot serve the server due to %s", err)
		}
	}()

	wait := gracefullyShutdown(context.Background(), 5*time.Second, map[string]operation{
		"http-server": func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
		"db-user": func(ctx context.Context) error {
			return db.Close()
		},
	})

	<-wait
}

func gracefullyShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {

		c := make(chan os.Signal, 1)

		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		// block until receive signal
		<-c

		log.Println("shutting down")

		// set timeout for the operation to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout)
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations async to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
