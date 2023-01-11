package pkg

import "errors"

var ErrDriverNotFound = errors.New("driver not found")
var ErrInvalidRequest = errors.New("invalid request")
var ErrInternalServer = errors.New("internal server error")
var ErrUnauthorized = errors.New("Unauthorized")
var ErrInvalidAPIKey = errors.New("invalid api key")
var ErrBulkCreateDriversLimit = errors.New("bulk create drivers limit exceeded")
