package errors

import (
	"net/http"

	"github.com/wisdom-oss/common-go/v2/types"
)

var MethodNotAllowed types.ServiceError = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#section-15.5.6",
	Status: http.StatusMethodNotAllowed,
	Title:  "Method Not Allowed",
	Detail: "The used HTTP method is not allowed on this route. Please check the documentation and your request",
}

var NotFound types.ServiceError = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#section-15.5.5",
	Status: http.StatusNotFound,
	Title:  "Route Not Found",
	Detail: "The requested path does not exist in this microservice. Please check the documentation and your request",
}
