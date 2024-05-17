package controller

import (
	"encoding/json"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/controller/serializer"
	dto "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/controller/serializer/input"
	"log/slog"
	"net/http"
)

type PaymentController struct {
	paymentUseCase contract.PaymentUseCase
	logger         *slog.Logger
	orderUseCase   contract.OrderUseCase
}

type GetLastPaymentStatus struct {
	OrderId       int
	PaymentStatus enum.PaymentStatus
}

func NewPaymentController(p contract.PaymentUseCase, logger *slog.Logger, orderUseCase contract.OrderUseCase) *PaymentController {
	return &PaymentController{
		paymentUseCase: p,
		logger:         logger,
		orderUseCase:   orderUseCase,
	}
}

// CreatePayment godoc
// @Summary      Start payment process
// @Description  Start payment process for a certain order
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        request   body      input.PaymentDto  true  "Payment properties"
// @Param        id   path      int  true  "Order ID"
// @Success      201  {object}  Response{data=string}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Failure      404  {object}  swagger.ResourceNotFoundResponse{data=interface{}}
// @Router       /orders/{id}/payments [post]
func (p *PaymentController) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var paymentDTO dto.PaymentDto
	var response Response

	if err := json.NewDecoder(r.Body).Decode(&paymentDTO); err != nil {
		p.logger.Error("Unable to decode the request body.  %v", slog.Any("error", err))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error decoding request body",
				Data:  nil,
			})
		return
	}

	if serialize := serializer.Validate(paymentDTO); len(serialize.Errors) > 0 {
		p.logger.Error("validate error", slog.Any("error", serialize))

		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(
			Response{
				Error: "Make sure all required fields are sent correctly",
				Data:  nil,
			})
		w.Write(jsonResponse)
		return
	}

	order := entity.Order{
		ID:          19,
		Amount:      1,
		StatusOrder: enum.RECEIVED,
	}

	qrCode, err := p.paymentUseCase.CreateQRCode(&order)

	if err != nil {
		p.logger.Error("error creating qr code", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse, _ := json.Marshal(
			Response{
				Error: "Error creating qr code",
				Data:  nil,
			})
		w.Write(jsonResponse)

		return
	}

	if qrCode == nil {
		response = Response{
			Error: "O pagamento para o pedido já foi efetuado",
			Data:  nil,
		}
	} else {
		response = Response{
			Error: "",
			Data:  qrCode,
		}
	}

	jsonResponse, _ := json.Marshal(response)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)

	return
}

// GetLastPaymentStatus godoc
// @Summary      Get status for last payment
// @Description  Get status for order last payment try
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  Response{data=string}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Router       /orders/{id}/status-payment [get]
func (p *PaymentController) GetLastPaymentStatus(w http.ResponseWriter, r *http.Request) {
	paymentIdParam := mux.Vars(r)["paymentId"]
	paymentId, err := strconv.Atoi(paymentIdParam)

	if err != nil {
		p.logger.Error("error to convert id order to int", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse, _ := json.Marshal(
			Response{
				Error: "Order id must be an integer",
				Data:  nil,
			})
		w.Write(jsonResponse)
		return
	}

	paymentStatus, err := p.paymentUseCase.GetLastPaymentStatus(paymentId)
	if err != nil {
		p.logger.Error("error getting last payment status", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse, _ := json.Marshal(
			Response{
				Error: "Error getting last payment status",
				Data:  nil,
			})
		w.Write(jsonResponse)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(
		Response{
			Error: "",
			Data: GetLastPaymentStatus{
				OrderId:       paymentId,
				PaymentStatus: paymentStatus,
			},
		})
	w.Write(jsonResponse)
	return
}

// Notification godoc
// @Summary      Payment confirmation webhook
// @Description  Payment confirmation webhook
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        request   body      input.PaymentDto  true  "Payment properties"
// @Param        id   path      int  true  "Order ID"
// @Success      201  {object}  interface{}
// @Failure      404  {object}  swagger.ResourceNotFoundResponse{data=interface{}}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Router       /orders/{id}/notification-payments [post]
func (p *PaymentController) Notification(w http.ResponseWriter, r *http.Request) {
	paymentIdParam := mux.Vars(r)["paymentId"]

	paymentId, err := strconv.Atoi(paymentIdParam)
	if err != nil {
		p.logger.Error("error to convert id order to int", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Order id must be an integer",
				Data:  nil,
			})

		return
	}

	if err = p.paymentUseCase.PaymentNotification(paymentId); err != nil {
		p.logger.Error("error processing payment notification", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error processing payment notification",
				Data:  nil,
			})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
