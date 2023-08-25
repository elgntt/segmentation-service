package main

import (
	"context"
	"log"

	"github.com/elgntt/avito-internship-2023/internal/api"
	"github.com/elgntt/avito-internship-2023/internal/config"
	"github.com/elgntt/avito-internship-2023/internal/pkg/db"
	repository "github.com/elgntt/avito-internship-2023/internal/repository/postgres"
	"github.com/elgntt/avito-internship-2023/internal/service"
)

func main() {
	dbConnectConfig, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := db.OpenDB(ctx, dbConnectConfig)
	if err != nil {
		log.Fatal(err)
	}

	r := api.New(
		service.New(
			repository.New(pool),
		),
	)

	log.Println("Server has been successfully started on the port :8080")
	log.Fatal(r.Run(":8080"))
}
