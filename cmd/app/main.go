package main

import (
	"context"
	"log"
	"time"

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

	service := service.New(
		repository.New(pool),
	)

	r := api.New(
		service,
	)

	go ClearExpiredSegmentsWorker(ctx, service)

	log.Println("Server has been successfully started on the port :8080")
	log.Fatal(r.Run(":8080"))
}

func ClearExpiredSegmentsWorker(ctx context.Context, s *service.Service) {
	for {
		workerInterval := time.NewTicker(1 * time.Minute)

		select {
		case <-workerInterval.C:
			err := s.DeleteExpiredUserSegments(ctx)
			if err != nil {
				log.Println("Worker err:", err)
			}
			log.Println("Success")
		}
	}

}
