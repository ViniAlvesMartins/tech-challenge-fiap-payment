package contract

import (
	response_payment_service "github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/modules/response/payment_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/enum"
)

type OrderUseCase interface {
	GetById(id int) (*entity.Order, error)
}

type PaymentUseCase interface {
	Create(payment *entity.Payment) error
	CreateQRCode(order *entity.Order) (*response_payment_service.CreateQRCode, error)
	GetLastPaymentStatus(orderId int) (enum.PaymentStatus, error)
	PaymentNotification(order *entity.Order) error
}

type CategoryUseCase interface {
	GetById(id int) (*entity.Category, error)
}

type ClientUseCase interface {
	GetClientByCpf(cpf int) (*entity.Client, error)
	GetClientById(id *int) (*entity.Client, error)
	Create(client entity.Client) (*entity.Client, error)
	GetAlreadyExists(cpf int, email string) (*entity.Client, error)
}

type ProductUseCase interface {
	Create(product entity.Product) (*entity.Product, error)
	Update(product entity.Product, id int) (*entity.Product, error)
	Delete(id int) error
	GetProductByCategory(categoryId int) ([]entity.Product, error)
	GetById(int int) (*entity.Product, error)
}
