package gateway

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Cinemaniac/config"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/routes"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/service"
	"github.com/saleh-ghazimoradi/Cinemaniac/slg"
	"github.com/saleh-ghazimoradi/Cinemaniac/utils"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Server() error {

	db, err := utils.DBConnection()
	if err != nil {
		return err
	}

	defer db.Close()

	server := &http.Server{
		Addr:         config.AppConfig.Server.Port,
		Handler:      routes.RegisterRoutes(db),
		IdleTimeout:  config.AppConfig.Server.IdleTimeout,
		ReadTimeout:  config.AppConfig.Server.ReadTimeout,
		WriteTimeout: config.AppConfig.Server.WriteTimeout,
		ErrorLog:     slog.NewLogLogger(slg.Logger.Handler(), slog.LevelError),
	}

	shutdownError := make(chan error)

	go func() {

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		slg.Logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err = server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		slg.Logger.Info("completing background tasks", "addr", server.Addr)

		service.WG.Wait()
		shutdownError <- nil
	}()

	slg.Logger.Info("starting server", "addr", server.Addr, "env", config.AppConfig.Server.Env)

	err = server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	slg.Logger.Info("stopped server", "addr", server.Addr)

	return nil
}
