package errors

import "github.com/wisdom-oss/common-go/v2/types"

var ErrMissingParameter = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "Request Missing Parameter",
	Detail: "The request is missing a required parameter. Check the error field for more information",
}

var ErrInvalidScope = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "Invalid Scope Set",
	Detail: "The request contained an invalid scope. Please check your request",
}

var ErrUnknownUser = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5",
	Status: 404,
	Title:  "Unknown User",
	Detail: "The user selected for this operation is not known",
}

var ErrRefreshTokenInvalid = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.4",
	Status: 403,
	Title:  "Invalid Refresh Token",
	Detail: "The refresh token is either expired or has been revoked",
}
