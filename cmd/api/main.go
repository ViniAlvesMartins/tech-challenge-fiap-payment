package main

import (
	"context"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-common/sns"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/use_case"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/config"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/handler/http_server"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/repository"
	mercadopago "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/service/external_payment/mercado_pago"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/service/order"
	snsproducer "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/service/sns"
	"github.com/go-resty/resty/v2"

	"github.com/ViniAlvesMartins/tech-challenge-fiap-common/dynamodb"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-common/uuid"
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

	db, err := dynamodb.NewConnection(ctx)
	if err != nil {
		logger.Error("error connecting to database", err)
		panic(err)
	}

	ordersHTTPClient := loadOrdersHttpClient(cfg.OrdersURL)

	orderService := order.NewService(ordersHTTPClient, logger)
	orderUseCase := use_case.NewOrderUseCase(orderService, logger)

	paymentRepository := repository.NewPaymentRepository(db, logger, loadUUID())
	snsConnection, err := sns.NewConnection(ctx, cfg.UpdateOrderStatusTopic)
	if err != nil {
		logger.Error("error connecting to sns", err)
		panic(err)
	}

	snsService := snsproducer.NewService(snsConnection)
	qrCodePaymentMethod := use_case.NewQRCode(mercadopago.NewPaymentGateway())
	paymentUseCase := use_case.NewPaymentUseCase(paymentRepository, qrCodePaymentMethod, snsService, logger)

	app := http_server.NewApp(logger, paymentUseCase, orderUseCase)
	app.Run()
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
