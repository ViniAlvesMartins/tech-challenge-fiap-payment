package controller

import (
	"encoding/json"
	"strconv"

	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"

	"github.com/gorilla/mux"

	"log/slog"
	"net/http"

	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/controller/serializer"
	dto "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/controller/serializer/input"
)

type PaymentController struct {
	paymentUseCase contract.PaymentUseCase
	logger         *slog.Logger
	orderUseCase   contract.OrderUseCase
}

type GetLastPaymentStatus struct {
	OrderId       int                `json:"id"`
	PaymentStatus enum.PaymentStatus `json:"status"`
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
// @Success      201  {object}  Response{data=string}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Failure      404  {object}  swagger.ResourceNotFoundResponse{data=interface{}}
// @Router       /payments [post]
func (p *PaymentController) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var paymentDTO dto.PaymentDto
	var response Response

	if err := json.NewDecoder(r.Body).Decode(&paymentDTO); err != nil {
		p.logger.Error("Unable to decode the request body.  %v", slog.Any("error", err))

		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse, _ := json.Marshal(
			Response{
				Error: "Error decoding request body",
				Data:  nil,
			})
		w.Write(jsonResponse)
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

	order, err := p.orderUseCase.GetById(paymentDTO.OrderId)

	if err != nil {
		p.logger.Error("error getting order", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(
			Response{
				Error: "Error getting order details",
				Data:  nil,
			})
		w.Write(jsonResponse)
		return
	}

	if order == nil {
		w.WriteHeader(http.StatusNotFound)
		jsonResponse, _ := json.Marshal(
			Response{
				Error: "Order not found",
				Data:  nil,
			})
		w.Write(jsonResponse)
		return
	}

	order.ID = paymentDTO.OrderId
	qrCode, err := p.paymentUseCase.CreateQRCode(r.Context(), order)

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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
	return
}

// GetLastPaymentStatus godoc
// @Summary      Get status for last payment
// @Description  Get status for order last payment try
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Payment ID"
// @Success      200  {object}  Response{data=string}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Router       /payments/{id}/status [get]
func (p *PaymentController) GetLastPaymentStatus(w http.ResponseWriter, r *http.Request) {
	paymentIdParam := mux.Vars(r)["id"]
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

	paymentStatus, err := p.paymentUseCase.GetLastPaymentStatus(r.Context(), paymentId)
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
// @Param        id   path      int  true  "Payment ID"
// @Success      201  {object}  interface{}
// @Failure      404  {object}  swagger.ResourceNotFoundResponse{data=interface{}}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Router       /payments/{id}/notification [post]
func (p *PaymentController) Notification(w http.ResponseWriter, r *http.Request) {
	paymentIdParam := mux.Vars(r)["id"]

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

	if err = p.paymentUseCase.ConfirmedPaymentNotification(r.Context(), paymentId); err != nil {
		p.logger.Error("error processing payment notification", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse, _ := json.Marshal(
			Response{
				Error: "Error processing payment notification",
				Data:  nil,
			})
		w.Write(jsonResponse)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Notification godoc
// @Summary      Payment cancelation endpoint
// @Description  Payment cancelation endpoint
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Payment ID"
// @Success      204  {object}  interface{}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Router       /payments/{id}/cancel [delete]
func (p *PaymentController) CancelPayment(w http.ResponseWriter, r *http.Request) {
	paymentIdParam := mux.Vars(r)["id"]

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

	if err = p.paymentUseCase.CanceledPaymentNotification(r.Context(), paymentId); err != nil {
		p.logger.Error("error canceling payment", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse, _ := json.Marshal(
			Response{
				Error: "error canceling payment",
				Data:  nil,
			})
		w.Write(jsonResponse)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
