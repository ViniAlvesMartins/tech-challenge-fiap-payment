package http_server

import (
	"context"
	"log/slog"
	"net/http"

	_ "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/doc/swagger"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/controller/controller"
	"github.com/swaggo/http-swagger/v2"

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
	router.HandleFunc("/payments/{orderId:[0-9]+}", paymentController.CreatePayment).Methods("POST")
	router.HandleFunc("/status-payment/{orderId:[0-9]+}", paymentController.GetLastPaymentStatus).Methods("GET")
	router.HandleFunc("/notification-payments/{orderId:[0-9]+}", paymentController.Notification).Methods("POST")

	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	return http.ListenAndServe(":8081", router)
}
