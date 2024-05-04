package use_case

import (
	"log/slog"

	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
)

type ProductUseCase struct {
	productRepository contract.ProductRepository
	logger            *slog.Logger
}

func NewProductUseCase(productRepository contract.ProductRepository, logger *slog.Logger) *ProductUseCase {
	return &ProductUseCase{
		productRepository: productRepository,
		logger:            logger,
	}
}

func (p *ProductUseCase) Create(product entity.Product) (*entity.Product, error) {

	product.Active = true
	productNew, err := p.productRepository.Create(product)

	if err != nil {
		return nil, err
	}

	return &productNew, nil
}

func (p *ProductUseCase) Update(product entity.Product, id int) (*entity.Product, error) {
	product.Active = true
	product.ID = id

	productUpdated, err := p.productRepository.Update(product)
	if err != nil {
		return nil, err
	}

	return productUpdated, nil
}

func (p *ProductUseCase) Delete(id int) error {
	err := p.productRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (p *ProductUseCase) GetProductByCategory(categoryId int) ([]entity.Product, error) {
	prod, err := p.productRepository.GetProductByCategory(categoryId)

	if err != nil {
		return nil, err
	}

	return prod, nil
}

func (p *ProductUseCase) GetById(id int) (*entity.Product, error) {
	prod, err := p.productRepository.GetById(id)

	if err != nil {
		return nil, err
	}

	return prod, nil
}
