package swagger

type ResourceNotFoundResponse struct {
	Error string      `json:"error" example:"Resource not found"`
	Data  interface{} `json:"data"`
}

type InternalServerErrorResponse struct {
	Error string      `json:"error" example:"Internal server error"`
	Data  interface{} `json:"data"`
}
