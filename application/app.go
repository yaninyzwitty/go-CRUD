package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router  http.Handler
	redisDb *redis.Client
}

func New() *App {
	app := &App{
		router:  loadRoutes(),
		redisDb: redis.NewClient(&redis.Options{}),
	}

	return app

}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}
	err := a.redisDb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	ch := make(chan error, 1)

	defer func() {
		if err := a.redisDb.Close(); err != nil {
			fmt.Println("failed to close redis client: ", err)
		}

	}()

	go func() {

		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
			close(ch)
		}

	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return server.Shutdown(timeout)

	}

}
