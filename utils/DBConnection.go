package utils

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Cinemaniac/config"
	"github.com/saleh-ghazimoradi/Cinemaniac/slg"
)

func postURI() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.AppConfig.Database.DatabaseHost, config.AppConfig.Database.DatabasePort, config.AppConfig.Database.DatabaseUser, config.AppConfig.Database.DatabasePassword, config.AppConfig.Database.DatabaseName, config.AppConfig.Database.DatabaseSSLMode)
}

func DBConnection() (*sql.DB, error) {
	uri := postURI()

	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.AppConfig.Database.MaxOpenConn)
	db.SetMaxIdleConns(config.AppConfig.Database.MaxIdleConn)
	db.SetConnMaxLifetime(config.AppConfig.Database.MaxLifetime)
	db.SetConnMaxIdleTime(config.AppConfig.Database.MaxIdleTime)

	slg.Logger.Info("Successfully connected to database")

	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig.CTX.Timeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
