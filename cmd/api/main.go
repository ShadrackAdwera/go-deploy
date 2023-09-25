package main

import (
	"context"
	"fmt"
	"go-k8s/internal/api"
	db "go-k8s/internal/db/sqlc"
	"go-k8s/internal/token"
	"go-k8s/internal/workers"
	"log"
	"os"
	"strings"

	"github.com/hibiken/asynq"
	zerolog "github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// err := godotenv.Load()

	// if err != nil {
	// 	fmt.Println(err)
	// 	log.Fatal("Error loading .env file")
	// }

	tokenKey := os.Getenv("TOKEN_KEY")
	pasetoMaker, err := token.NewPasetoMaker(strings.TrimSpace(tokenKey))

	if err != nil {
		fmt.Println(err)
		log.Fatal("Error creating the token maker")
	}

	dbUrl := os.Getenv("DB_URL")

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error creating connection pool")
	}

	migrationUrl := os.Getenv("MIGRATION_URL")

	// run DB migrations
	runDBMigration(migrationUrl, dbUrl)

	store := db.NewStore(pool)

	redisAddress := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	clientOpts := asynq.RedisClientOpt{
		Addr:     redisAddress,
		Password: redisPassword,
	}

	distro := workers.NewRedisTaskDistributor(clientOpts)

	srv := api.NewServer(store, pasetoMaker, distro)

	serverAddress := os.Getenv("SERVER_ADDRESS")
	go startTaskProcessor(clientOpts, store)
	if err := srv.StartServer(serverAddress); err != nil {
		fmt.Println(err)
		log.Fatal("Error starting the HTTP server")
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		fmt.Println(err)
		log.Fatal("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Println(err)
		log.Fatal("failed to run migrate up")
	}

	fmt.Println("db migrated successfully")
}

func startTaskProcessor(opts asynq.RedisClientOpt, store db.TxStore) {
	processor := workers.NewRedisTaskProcessor(opts, store)

	err := processor.Start()
	zerolog.Info().Msg("connecting to REDIS processor . . . ")

	if err != nil {
		zerolog.Fatal().Err(err).Msg("error starting the redis task processor")
		return
	}
	zerolog.Info().Msg("redis task processor started . . . ")
}
