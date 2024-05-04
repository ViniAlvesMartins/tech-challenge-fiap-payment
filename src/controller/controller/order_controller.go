package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/controller/serializer/input"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/controller/serializer/output"

	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/enum"

	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"

	"github.com/gorilla/mux"
)

type OrderController struct {
	orderUseCase   contract.OrderUseCase
	productUseCase contract.ProductUseCase
	clientUseCase  contract.ClientUseCase
	logger         *slog.Logger
}

func NewOrderController(orderUseCase contract.OrderUseCase, productUseCase contract.ProductUseCase, clientUseCase contract.ClientUseCase, logger *slog.Logger) *OrderController {
	return &OrderController{
		orderUseCase:   orderUseCase,
		productUseCase: productUseCase,
		clientUseCase:  clientUseCase,
		logger:         logger,
	}
}

// CreateOrder godoc
// @Summary      Create order
// @Description  Place a new order
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        request   body      input.OrderDto  true  "Order properties"
// @Success      201  {object}  Response{data=output.OrderDto}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Failure      404  {object}  swagger.ResourceNotFoundResponse{data=interface{}}
// @Router       /orders [post]
func (o *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderDto input.OrderDto

	if err := json.NewDecoder(r.Body).Decode(&orderDto); err != nil {
		o.logger.Error("unable to decode the request body", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Unable to decode the request body",
				Data:  nil,
			})
		return
	}

	if prods := len(orderDto.Products); prods < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Error: "Product is required",
			Data:  nil,
		})
		return
	}

	var products []*entity.Product
	for _, p := range orderDto.Products {
		product, err := o.productUseCase.GetById(p.ID)

		if err != nil {
			o.logger.Error("error getting product by id", slog.String("message", err.Error()))

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{
				Error: "Error finding product",
				Data:  nil,
			})
			return
		}

		if product == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(
				Response{
					Error: fmt.Sprintf("Product of id %d not found", p.ID),
					Data:  nil,
				})
			return
		}

		products = append(products, product)
	}

	client, err := o.clientUseCase.GetClientById(orderDto.ClientId)
	if client == nil && orderDto.ClientId != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Client not found",
				Data:  nil,
			})
		return
	}

	order, err := o.orderUseCase.Create(orderDto.ConvertToEntity(), products)
	if err != nil {
		o.logger.Error("error creating order", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error creating order",
				Data:  nil,
			})
		return
	}

	orderOutput := output.OrderFromEntity(*order)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(
		Response{
			Error: "",
			Data:  orderOutput,
		})
}

// FindOrders godoc
// @Summary      List orders
// @Description  List orders by status
// @Tags         Orders
// @Produce      json
// @Success      200  {object}  Response{data=[]output.OrderDto}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Router       /orders [get]
func (o *OrderController) FindOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := o.orderUseCase.GetAll()
	if err != nil {
		o.logger.Error("error listing orders", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error listing orders",
				Data:  nil,
			})
		return
	}

	ordersOutput := output.OrderListFromEntity(*orders)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Error: "",
		Data:  ordersOutput,
	})
}

// GetOrderById godoc
// @Summary      Find order
// @Description  Find order by id
// @Tags         Orders
// @Produce      json
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  Response{data=output.OrderDto}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Failure      404  {object}  swagger.ResourceNotFoundResponse{data=interface{}}
// @Router       /orders/{id} [get]
func (o *OrderController) GetOrderById(w http.ResponseWriter, r *http.Request) {
	orderIdParam := mux.Vars(r)["orderId"]

	id, err := strconv.Atoi(orderIdParam)
	if err != nil {
		o.logger.Error("error to convert id order to int", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Order id must be an integer",
				Data:  nil,
			})
		return
	}

	order, err := o.orderUseCase.GetById(id)
	if err != nil {
		o.logger.Error("error finding order", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error finding order",
				Data:  nil,
			})
		return
	}

	if order == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Order not found",
				Data:  nil,
			})
		return
	}

	orderOutput := output.OrderFromEntity(*order)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Error: "",
		Data:  orderOutput,
	})
}

// UpdateOrderStatusById godoc
// @Summary      Find order
// @Description  Find order by id
// @Tags         Orders
// @Produce      json
// @Param        id   path      int  true  "Order ID"
// @Param        request   body      input.StatusOrderDto  true  "Order status"
// @Success      204  {object}  interface{}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Failure      404  {object}  swagger.ResourceNotFoundResponse{data=interface{}}
// @Router       /orders/{id} [patch]
func (o *OrderController) UpdateOrderStatusById(w http.ResponseWriter, r *http.Request) {
	orderIdParam := mux.Vars(r)["orderId"]
	orderId, err := strconv.Atoi(orderIdParam)

	if err != nil {
		o.logger.Error("error converting orderId to int", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Error: "Error to convert orderId to int",
			Data:  nil,
		})
		return
	}

	var statusOrderDto input.StatusOrderDto
	if err := json.NewDecoder(r.Body).Decode(&statusOrderDto); err != nil {
		o.logger.Error("unable to decode the request body", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Unable to decode the request body",
				Data:  nil,
			})
		return
	}

	if !enum.ValidateStatus(statusOrderDto.Status) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Error: "Invalid status",
			Data:  nil,
		})
		return
	}

	order, err := o.orderUseCase.GetById(orderId)
	if err != nil {
		o.logger.Error("error getting order by id", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Error: "Error to get order",
			Data:  nil,
		})
		return
	}

	if order == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Order not found",
				Data:  nil,
			})
		return
	}

	if err := o.orderUseCase.UpdateStatusById(orderId, enum.StatusOrder(statusOrderDto.Status)); err != nil {
		o.logger.Error("error updating status by id", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Error: "Error updating status",
			Data:  nil,
		})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
