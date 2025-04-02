package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/saswatax/rss-aggregator/internal/config"
	"github.com/saswatax/rss-aggregator/internal/database"
	"github.com/saswatax/rss-aggregator/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Falied to load config:", err)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.Env.DB_URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close(ctx)
	db := database.New(conn)

	s := server.Server{
		Config: cfg,
		DB:     db,
	}

	srv := server.NewServer(s)

	log.Printf("Server starting on port %v", os.Getenv("PORT"))
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
