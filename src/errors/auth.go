package errors

import (
	"errors"
	"github.com/ACLzz/keystore-server/src/utils"
)

var PasswordMaxLengthError = errors.New(lengthMax("password", utils.PasswordMaxLengthLimit))
var PasswordMinLengthError = errors.New(lengthMin("password", utils.PasswordMinLengthLimit))
var LoginMaxLengthError = errors.New(lengthMax("login", utils.LoginMaxLengthLimit))
var LoginMinLengthError = errors.New(lengthMin("login", utils.LoginMinLengthLimit))

var PasswordLocaleError = errors.New(localeError("password"))
var LoginLocaleError = errors.New(localeError("login"))
var TokenDeniedSymbolsError = errors.New("your token includes denied symbols")

var InvalidCredentials = errors.New("invalid credentials")

var NoLoginError = errors.New(noField("login"))
var NoPasswordError = errors.New(noField("password"))

var NoToken = errors.New(noField("token"))
var InvalidToken = errors.New("invalid token")
var ExpiredToken = errors.New("your token has expired")

var UserExists = errors.New("user already exists")
var UserNotExists = errors.New("user not exists")

var RegistrationDisabled = errors.New("administrator disabled registration for new users")
