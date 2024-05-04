package input

import "github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"

type ClientDto struct {
	Cpf   int    `json:"cpf" validate:"required" error:"Campo cpf é obrigatorio"`
	Name  string `json:"name" validate:"required" error:"Campo nome é obrigatorio"`
	Email string `json:"email" validate:"required" error:"Campo email é obrigatorio"`
}

func (c *ClientDto) ConvertEntity() entity.Client {
	return entity.Client{
		Name:  c.Name,
		Cpf:   c.Cpf,
		Email: c.Email,
	}
}
