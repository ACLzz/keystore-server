package errors

import (
	"errors"
	"fmt"
)

// Internal
var InternalError = errors.New("internal error")

// Requests
var EmptyBodyError = errors.New("empty body")
var InvalidBodyError = errors.New("invalid body")

// Auth
var PasswordMaxLengthError = errors.New(lengthMax("password", PasswordMaxLengthLimit))
var LoginMaxLengthError = errors.New(lengthMax("login", LoginMaxLengthLimit))
var PasswordMinLengthError = errors.New(lengthMin("password", PasswordMinLengthLimit))
var LoginMinLengthError = errors.New(lengthMin("login", LoginMinLengthLimit))

var PasswordLocaleError = errors.New("password must contain only english or digit symbols")
var LoginLocaleError = errors.New("login must contain only english or digit symbols")
var TokenDeniedSymbolsError = errors.New("your token includes denied symbols")

var InvalidCredentials = errors.New("invalid credentials")

var NoLoginError = errors.New(noField("login"))
var NoPasswordError = errors.New(noField("password"))

var NoToken = errors.New(noField("token"))
var InvalidToken = errors.New("invalid token")
var ExpiredToken = errors.New("your token has expired")

// database.User
var UserExists = errors.New("user already exists")
var UserNotExists = errors.New("user not exists")


func noField(field string) string {
	return fmt.Sprintf("no \"%s\" was sent", field)
}

func lengthMax(field string, limit int) string {
	return fmt.Sprintf("%s can't be more than %d symbols", field, limit)
}

func lengthMin(field string, limit int) string {
	return fmt.Sprintf("%s must be at least %d symbols", field, limit)
}
