package entity

type Product struct {
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	NameProduct string  `json:"name_product"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	CategoryId  int     `json:"category_id"`
	Active      bool    `json:"active"`
}
