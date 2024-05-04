package entity

type Client struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Cpf   int    `json:"cpf"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
