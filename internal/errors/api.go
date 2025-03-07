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

var ErrUserDisabled = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.4",
	Status: 403,
	Title:  "User Disabled",
	Detail: "The user has been disabled.",
}

var ErrRefreshTokenInvalid = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.4",
	Status: 401,
	Title:  "Invalid Refresh Token",
	Detail: "The refresh token is either expired or has been revoked",
}

var ErrBadService = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "Unknown Service",
	Detail: "The service provided in the request is unknown",
}

var ErrUnknownService = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5",
	Status: 404,
	Title:  "Unknown Service",
	Detail: "The service provided in the request is unknown",
}

var ErrInvalidClientCredentials = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.2",
	Status: 401,
	Title:  "Invalid Client Credentials",
	Detail: "The supplied client credentials are not valid",
}

var ErrInvalidClientScopeRequested = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.2",
	Status: 400,
	Title:  "Unsupported Client Scope Requested",
	Detail: "At least one of the requested scope is not supported for client credentials",
}

var ErrPermissionMismatch = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.2",
	Status: 400,
	Title:  "Permission Mismatch",
	Detail: "Assigning scopes outside of the permission scope of the current user is not supported",
}

var ErrInvalidClientID = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.2",
	Status: 400,
	Title:  "Invalid Client ID Format",
	Detail: "Invalid Client ID provided. Please ensure you used an UUIDv4",
}
