package gateway

import (
	"github.com/saleh-ghazimoradi/Cinemaniac/config"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/routes"
	"github.com/saleh-ghazimoradi/Cinemaniac/slg"
	"log/slog"
	"net/http"
	"os"
)

func Server() error {
	server := &http.Server{
		Addr:         config.AppConfig.Server.Port,
		Handler:      routes.RegisterRoutes(),
		IdleTimeout:  config.AppConfig.Server.IdleTimeout,
		ReadTimeout:  config.AppConfig.Server.ReadTimeout,
		WriteTimeout: config.AppConfig.Server.WriteTimeout,
		ErrorLog:     slog.NewLogLogger(slg.Logger.Handler(), slog.LevelError),
	}

	slg.Logger.Info("starting server", "addr", server.Addr, "env", config.AppConfig.Server.Env)

	if err := server.ListenAndServe(); err != nil {
		slg.Logger.Error("error starting server", "error", err)
		os.Exit(1)
	}

	return nil
}
