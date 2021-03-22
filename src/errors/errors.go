package errors

import (
	"errors"
	"fmt"
)

// Internal
var InternalError = errors.New("internal error")

// Requests
var EmptyBodyError = errors.New("empty body")
var InvalidBodyError = errors.New("invalid boy")

// Auth
var PasswordMaxLengthError = errors.New(fmt.Sprintf("password can't be more than %d symbols", PasswordMaxLengthLimit))
var LoginMaxLengthError = errors.New(fmt.Sprintf("username can't be more than %d symbols", UsernameMaxLengthLimit))
var PasswordMinLengthError = errors.New(fmt.Sprintf("password must be at least %d symbols", PasswordMinLengthLimit))
var LoginMinLengthError = errors.New(fmt.Sprintf("username must be at least %d symbols", UsernameMinLengthLimit))

var PasswordLocaleError = errors.New("password must contain only english or digit symbols")
var LoginLocaleError = errors.New("login must contain only english or digit symbols")

var NoLoginError = errors.New("no \"login\" in body provided")
var NoPasswordError = errors.New("no \"password\" in body provided")

// database.User
var UserExists = errors.New("user already exists")
var UserNotExists = errors.New("user not exists")
