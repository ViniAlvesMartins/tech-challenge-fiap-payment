package main

import (
	"context"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/infra"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/use_case"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/external/database/dynamodb"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/external/handler/http_server"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/external/repository"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/external/service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/pkg/uuid"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"os"
)

// @title           Ze Burguer Payment APIs
// @version         1.0
func main() {
	var err error
	var ctx = context.Background()
	var logger = loadLogger()

	cfg, err := loadConfig()

	if err != nil {
		logger.Error("error loading config", err)
		panic(err)
	}

	db, err := dynamodb.NewConnection(cfg)

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

func loadConfig() (infra.Config, error) {
	return infra.NewConfig()
}

func loadLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}

func loadOrdersHttpClient(ordersURL string) *resty.Client {
	return resty.New().SetBaseURL(ordersURL)
}
