package response

import "github.com/go-playground/validator/v10"

const (
	StatusOk    = "ok"
	StatusError = "error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func NewOkResponse() Response {
	return Response{Status: StatusOk}
}

func NewErrorResponse(err error) Response {
	return Response{Status: StatusError, Error: err.Error()}
}

func NewValidationErrorResponse(errors validator.ValidationErrors) *Response {
	var err string

	for _, e := range errors {
		err += e.Error() + "\n"
	}

	return &Response{Status: StatusError, Error: err}
}
