package controller

import (
	"encoding/json"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/controller/serializer/output"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/controller/serializer"

	dto "github.com/ViniAlvesMartins/tech-challenge-fiap/src/controller/serializer/input"
)

type ClientController struct {
	clientUseCase contract.ClientUseCase
	logger        *slog.Logger
}

func NewClientController(clientUseCase contract.ClientUseCase, logger *slog.Logger) *ClientController {
	return &ClientController{
		clientUseCase: clientUseCase,
		logger:        logger,
	}
}

// CreateClient godoc
// @Summary      Create client
// @Description  Add new client
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        request   body      input.ClientDto  true  "Client properties"
// @Success      201  {object}  Response{error=string,data=output.ClientDto}
// @Failure      500  {object}  swagger.ResourceNotFoundResponse{data=interface{}}
// @Router       /clients [post]
func (c *ClientController) CreateClient(w http.ResponseWriter, r *http.Request) {
	var clientDto dto.ClientDto

	if err := json.NewDecoder(r.Body).Decode(&clientDto); err != nil {
		c.logger.Error("unable to decode the request body", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Unable to decode the request body",
				Data:  nil,
			})
		return
	}

	validate := serializer.Validate(clientDto)
	if len(validate.Errors) > 0 {
		c.logger.Error("validate error", slog.Any("error", validate))

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Invalid body, make sure all required fields are sent",
				Data:  nil,
			})
		return
	}

	validClient, err := c.clientUseCase.GetAlreadyExists(clientDto.Cpf, clientDto.Email)
	if err != nil {
		c.logger.Error("error validating client", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error validating client",
				Data:  nil,
			})
		return
	}

	if validClient != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Client already exists",
				Data:  nil,
			})
		return
	}

	client, err := c.clientUseCase.Create(clientDto.ConvertEntity())
	if err != nil {
		c.logger.Error("error creating client", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: err.Error(),
				Data:  nil,
			})
		return
	}

	clientOutput := output.ClientFromEntity(*client)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(
		Response{
			Error: "",
			Data:  clientOutput,
		})
}

// GetClientByCpf godoc
// @Summary      Show client details
// @Description  Get client by cpf
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        cpf   query      integer  true  "Client cpf"
// @Success      200  {object}  Response{data=output.ClientDto}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Failure      404  {object}  swagger.ResourceNotFoundResponse{data=interface{}}
// @Router       /clients [get]
func (c *ClientController) GetClientByCpf(w http.ResponseWriter, r *http.Request) {
	cpfParam := r.URL.Query().Get("cpf")

	cpf, err := strconv.Atoi(cpfParam)
	if err != nil {
		c.logger.Error("error to convert cpf to int", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Make sure document is an int",
				Data:  nil,
			})
		return
	}

	client, err := c.clientUseCase.GetClientByCpf(cpf)
	if err != nil {
		c.logger.Error("error finding client", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error finding client",
				Data:  nil,
			})
		return
	}

	if client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Client not found",
				Data:  nil,
			})
		return
	}

	clientOutput := output.ClientFromEntity(*client)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		Response{
			Error: "",
			Data:  clientOutput,
		})
}
