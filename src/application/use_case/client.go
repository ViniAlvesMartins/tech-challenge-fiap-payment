package use_case

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"

	"log/slog"
)

type ClientUseCase struct {
	clientRepository contract.ClientRepository
	logger           *slog.Logger
}

func NewClientUseCase(clientRepository contract.ClientRepository, logger *slog.Logger) *ClientUseCase {
	return &ClientUseCase{
		clientRepository: clientRepository,
		logger:           logger,
	}
}

func (c *ClientUseCase) Create(client entity.Client) (*entity.Client, error) {

	clientNew, err := c.clientRepository.Create(client)

	if err != nil {
		return nil, err
	}

	return &clientNew, nil
}

func (c *ClientUseCase) GetClientByCpf(cpf int) (*entity.Client, error) {
	client, err := c.clientRepository.GetClientByCpf(cpf)

	if err != nil {
		return nil, err
	}

	return client, nil

}

func (c *ClientUseCase) GetClientById(id *int) (*entity.Client, error) {
	client, err := c.clientRepository.GetClientById(id)

	if err != nil {
		return nil, err
	}

	return client, nil

}

func (c *ClientUseCase) GetAlreadyExists(cpf int, email string) (*entity.Client, error) {
	client, err := c.clientRepository.GetAlreadyExists(cpf, email)

	if err != nil {
		return nil, err
	}

	return client, nil

}
