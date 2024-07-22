package http_server

import (
	"context"
	"errors"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/doc/swagger"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func (e *App) Run() {
	router := mux.NewRouter()

	paymentController := controller.NewPaymentController(e.paymentUseCase, e.logger, e.orderUseCase)
	router.HandleFunc("/payments", paymentController.CreatePayment).Methods("POST")
	router.HandleFunc("/payments/{id:[0-9]+}/status", paymentController.GetLastPaymentStatus).Methods("GET")
	router.HandleFunc("/payments/{id:[0-9]+}/cancel", paymentController.CancelPayment).Methods("DELETE")
	router.HandleFunc("/payments/{id:[0-9]+}/notification", paymentController.Notification).Methods("POST")
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	swagger.SwaggerInfo.Title = "Ze Burguer Payment API"
	swagger.SwaggerInfo.Version = "1.0"

	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatal(err)
	}
}
