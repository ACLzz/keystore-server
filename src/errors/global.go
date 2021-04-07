package errors

import "errors"

// Internal
var InternalError = errors.New("internal error")

// Requests
var EmptyBodyError = errors.New("empty body")
var InvalidBodyError = errors.New("invalid body")
