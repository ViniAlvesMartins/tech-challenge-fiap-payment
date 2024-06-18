package entity

type Client struct {
	ID    int    `json:"id"`
	Cpf   int    `json:"cpf"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
