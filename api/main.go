package main

import (
	"context"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/use_case"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/config"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/handler/http_server"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/repository"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/pkg/dynamodb"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/pkg/uuid"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"os"
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

	db, err := dynamodb.NewConnection(ctx, cfg.DynamoDBRegion, cfg.DynamoDBUrl, cfg.DynamoDBAccessKey, cfg.DynamoDBSecretAccess)
	if err != nil {
		logger.Error("error connecting to database", err)
		panic(err)
	}

	ordersHTTPClient := loadOrdersHttpClient(cfg.OrdersURL)

	orderService := service.NewOrderService(ordersHTTPClient, logger)
	orderUseCase := use_case.NewOrderUseCase(orderService, logger)

	paymentRepository := repository.NewPaymentRepository(db, logger, loadUUID())
	snsService := service.NewSnsService()
	externalPaymentService := service.NewExternalPayment()
	paymentUseCase := use_case.NewPaymentUseCase(paymentRepository, externalPaymentService, snsService, logger)

	app := http_server.NewApp(logger, paymentUseCase, orderUseCase)

	err = app.Run(ctx)

	if err != nil {
		logger.Error("error running application", err)
		panic(err)
	}
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

func loadOrdersHttpClient(ordersURL string) *resty.Client {
	return resty.New().SetBaseURL(ordersURL)
}
