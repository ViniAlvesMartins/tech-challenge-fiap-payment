package output

type ClientDto struct {
	ID    int    `json:"id"`
	Cpf   int    `json:"cpf"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
