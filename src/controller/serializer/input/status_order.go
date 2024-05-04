package input

type StatusOrderDto struct {
	Status string `json:"status" validate:"required"`
}
