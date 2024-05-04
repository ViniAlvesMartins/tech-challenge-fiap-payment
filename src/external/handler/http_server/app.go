package http_server

import (
	"context"
	"log/slog"
	"net/http"

	_ "github.com/ViniAlvesMartins/tech-challenge-fiap/doc/swagger"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/controller/controller"
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
	router.HandleFunc("/payments", paymentController.CreatePayment).Methods("POST")
	router.HandleFunc("/status-payment", paymentController.GetLastPaymentStatus).Methods("GET")
	router.HandleFunc("/notification-payments", paymentController.Notification).Methods("POST")

	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	return http.ListenAndServe(":8081", router)
}
