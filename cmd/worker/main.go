package main

import (
	"context"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-common/dynamodb"
	sqsservice "github.com/ViniAlvesMartins/tech-challenge-fiap-common/sqs"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-common/uuid"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/use_case"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/config"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/repository"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/service/sqs"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var err error
	var ctx, cancel = context.WithCancel(context.Background())
	var logger = loadLogger()

	fmt.Println("Initializing worker...")
	cfg, err := loadConfig()

	if err != nil {
		logger.Error("error loading config", err)
		panic(err)
	}

	db, err := dynamodb.NewConnection(ctx)
	if err != nil {
		logger.Error("error connecting to database", err)
		panic(err)
	}

	paymentRepository := repository.NewPaymentRepository(db, logger, loadUUID())

	consumer, err := sqsservice.NewConnection(ctx, cfg.ProductionFailedQueue, 1, 20)
	if err != nil {
		logger.Error("error connecting to sqs", err)
		panic(err)
	}

	paymentUseCase := use_case.NewPaymentUseCase(paymentRepository, nil, nil, logger)
	failedProductionHandler := sqs.NewFailedProductionHandler(paymentUseCase)

	fmt.Println("Starting consumer...")
	failedProductionConsumer := sqs.NewConsumer(consumer, failedProductionHandler)

	var wg sync.WaitGroup
	wg.Add(1)
	go failedProductionConsumer.Start(ctx, &wg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	cancel()
	wg.Wait()
	fmt.Println("Finishing worker...")
}

func loadUUID() uuid.Interface {
	return &uuid.UUID{}
}

func loadConfig() (config.Config, error) {
	return config.NewConfig()
}

func loadLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}
