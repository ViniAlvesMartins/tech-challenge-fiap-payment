package controller

import (
	"encoding/json"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/controller/serializer"
	dto "github.com/ViniAlvesMartins/tech-challenge-fiap/src/controller/serializer/input"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/controller/serializer/output"

	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductController struct {
	productUseCase  contract.ProductUseCase
	categoryUseCase contract.CategoryUseCase
	logger          *slog.Logger
}

func NewProductController(productUseCase contract.ProductUseCase, categoryUseCase contract.CategoryUseCase, logger *slog.Logger) *ProductController {
	return &ProductController{
		productUseCase:  productUseCase,
		logger:          logger,
		categoryUseCase: categoryUseCase,
	}
}

// CreateProduct godoc
// @Summary      Create product
// @Description  Place a new product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        request   body      input.ProductDto  true  "Product properties"
// @Success      200  {object}  Response{error=string,data=output.ProductDto}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Router       /products [post]
func (p *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productDto dto.ProductDto

	if err := json.NewDecoder(r.Body).Decode(&productDto); err != nil {
		p.logger.Error("Unable to decode the request body.  %v", slog.Any("error", err))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error decoding request body",
				Data:  nil,
			})
		return
	}

	if serialize := serializer.Validate(productDto); len(serialize.Errors) > 0 {
		p.logger.Error("validate error", slog.Any("error", serialize))

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Make all required fields are sent correctly",
				Data:  nil,
			})
		return
	}

	category, err := p.categoryUseCase.GetById(productDto.CategoryId)
	if err != nil {
		p.logger.Error("validate getting category by id", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error getting category",
				Data:  nil,
			})
		return
	}

	if category == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Category not found",
				Data:  nil,
			})
		return
	}

	product, err := p.productUseCase.Create(productDto.ConvertToEntity())
	if err != nil {
		p.logger.Error("error creating product", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error creating product",
				Data:  nil,
			})
		return
	}

	productOutput := output.ProductFromEntity(*product)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(
		Response{
			Error: "",
			Data:  productOutput,
		})
}

// UpdateProduct godoc
// @Summary      Update product
// @Description  Update product properties
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Param        request   body      input.ProductDto  true  "Product properties"
// @Success      200  {object}  Response{error=string,data=output.ProductDto}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Router       /products/{id} [put]
func (p *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var productDto dto.ProductDto

	productIdParam, ok := mux.Vars(r)["productId"]
	if !ok {
		p.logger.Error("id is missing in parameters")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Id is missing in parameters",
				Data:  nil,
			})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&productDto); err != nil {
		p.logger.Error("Unable to decode the request body.  %v", slog.Any("error", err))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error decoding request body",
				Data:  nil,
			})
		return
	}

	productId, err := strconv.Atoi(productIdParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Id is not a number",
				Data:  nil,
			})
		return
	}

	if serialize := serializer.Validate(productDto); len(serialize.Errors) > 0 {
		p.logger.Error("validate error", slog.Any("error", serialize))

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Make sure all required fields are sent correctly",
				Data:  nil,
			})
		return
	}

	validateProduct, err := p.productUseCase.GetById(productId)
	if err != nil {
		p.logger.Error("error getting product by id", slog.Any("error", err))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error finding product",
				Data:  nil,
			})
		return
	}

	if validateProduct == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Product not found",
				Data:  nil,
			})
		return
	}

	productDomain := productDto.ConvertToEntity()
	product, err := p.productUseCase.Update(productDomain, productId)
	if err != nil {
		p.logger.Error("error updating product data", slog.Any("error", err))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error updating product data",
				Data:  nil,
			})
		return
	}

	productOutput := output.ProductFromEntity(*product)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		Response{
			Error: "",
			Data:  productOutput,
		})
}

// DeleteProduct godoc
// @Summary      Delete product
// @Description  Remove product from list
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  interface{}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Router       /products/{id} [delete]
func (p *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productIdParam, ok := mux.Vars(r)["productId"]
	if !ok {
		p.logger.Error("id is missing in parameters")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Id is missing in parameters",
				Data:  nil,
			})
		return
	}

	productId, err := strconv.Atoi(productIdParam)
	if err != nil {
		p.logger.Error("Error to convert productId to int.  %v", err)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Id is not a number",
				Data:  nil,
			})
		return
	}

	validateProduct, err := p.productUseCase.GetById(productId)
	if err != nil {
		p.logger.Error("error getting product by id", slog.Any("error", err))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error finding product",
				Data:  nil,
			})
		return
	}

	if validateProduct == nil || validateProduct.Active == false {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Product not found",
				Data:  nil,
			})
		return
	}

	if err := p.productUseCase.Delete(productId); err != nil {
		p.logger.Error("error deleting product", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error deleting product",
				Data:  nil,
			})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetProductByCategory godoc
// @Summary      List product by category
// @Description  List products from a certain category
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      200  {object}  Response{data=[]output.ProductDto}
// @Failure      500  {object}  swagger.InternalServerErrorResponse{data=interface{}}
// @Router       /categories/{id}/products [get]
func (p *ProductController) GetProductByCategory(w http.ResponseWriter, r *http.Request) {
	categoryIdParam := mux.Vars(r)["categoryId"]

	categoryId, err := strconv.Atoi(categoryIdParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Id is missing in parameters",
				Data:  nil,
			})
		return
	}

	category, err := p.categoryUseCase.GetById(categoryId)
	if err != nil {
		p.logger.Error("error getting category by id", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error finding category",
				Data:  nil,
			})
		return
	}

	if category == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Category not found",
				Data:  nil,
			})
		return
	}

	products, err := p.productUseCase.GetProductByCategory(categoryId)
	if err != nil {
		p.logger.Error("error getting products by category", slog.Any("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Error finding products",
				Data:  nil,
			})
		return
	}

	if len(products) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			Response{
				Error: "Product not found",
				Data:  nil,
			})
		return
	}

	productsOutput := output.ProductListFromEntity(products)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		Response{
			Error: "",
			Data:  productsOutput,
		})
}
