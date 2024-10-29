package main

import (
	"context"
	"flag"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"

	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
	"github.com/ArturSaga/chat-server/internal/app"
)

type server struct {
	desc.UnimplementedChatApiServer
	pool *pgxpool.Pool
}

const timeout = 15

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
