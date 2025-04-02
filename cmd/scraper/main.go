package scraper

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/saswatax/rss-aggregator/internal/config"
	"github.com/saswatax/rss-aggregator/internal/database"
	"github.com/saswatax/rss-aggregator/internal/scraper"
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

	scraper.StartScraping(db, 10, 10*time.Minute)
}
