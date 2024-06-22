package http_server

import (
	"context"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/doc/swagger"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log/slog"
	"net/http"

	_ "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/doc/swagger"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/controller/controller"
	"github.com/gorilla/mux"
)

type App struct {
	logger         *slog.Logger
	paymentUseCase contract.PaymentUseCase
	orderUseCase   contract.OrderUseCase
}

func NewApp(logger *slog.Logger,
	paymentUseCase contract.PaymentUseCase,
	orderUseCase contract.OrderUseCase,
) *App {
	return &App{
		logger:         logger,
		paymentUseCase: paymentUseCase,
		orderUseCase:   orderUseCase,
	}
}

func (e *App) Run(ctx context.Context) error {
	router := mux.NewRouter()

	paymentController := controller.NewPaymentController(e.paymentUseCase, e.logger, e.orderUseCase)
	router.HandleFunc("/payments", paymentController.CreatePayment).Methods("POST")
	router.HandleFunc("/payments/{paymentId:[0-9]+}/status", paymentController.GetLastPaymentStatus).Methods("GET")
	router.HandleFunc("/payments/{paymentId:[0-9]+}/notification-payments", paymentController.Notification).Methods("POST")

	swagger.SwaggerInfo.Title = "Ze Burguer Payment API"
	swagger.SwaggerInfo.Version = "1.0"

	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	return http.ListenAndServe(":8081", router)
}
