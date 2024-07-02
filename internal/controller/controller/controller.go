package controller

type Response struct {
	Error string      `json:"error" example:""`
	Data  interface{} `json:"data"`
}
