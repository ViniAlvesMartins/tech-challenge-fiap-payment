package use_case

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"

	"log/slog"
)

type CategoryUseCase struct {
	categoryRepository contract.CategoryRepository
	logger             *slog.Logger
}

func NewCategoryUseCase(categoryRepository contract.CategoryRepository, logger *slog.Logger) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepository: categoryRepository,
		logger:             logger,
	}
}

func (c *CategoryUseCase) GetById(id int) (*entity.Category, error) {

	category, err := c.categoryRepository.GetById(id)

	if err != nil {
		return nil, err
	}

	return category, nil
}
