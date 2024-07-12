package main

import (
	"context"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/use_case"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/config"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/repository"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/service/sqs"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/pkg/dynamodb"
	sqsservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/pkg/sqs"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/pkg/uuid"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var err error
	var ctx = context.Background()
	var logger = loadLogger()

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

	failedProductionConsumer := sqs.NewConsumer(consumer, failedProductionHandler)
	failedProductionConsumer.Start(ctx)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
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
