package errors

import (
	"errors"
	"github.com/ACLzz/keystore-server/src/utils"
)

var PasswordNotExist = errors.New("password with that id does not exist")
var InvalidPasswordId = errors.New("you sent invalid password id")

var PasswordTitleMaxLengthError = errors.New(lengthMax("title", utils.PasswordTitleMaxLengthLimit))
var PasswordTitleMinLengthError = errors.New(lengthMin("title", utils.PasswordTitleMinLengthLimit))
var PasswordLoginMaxLengthError = errors.New(lengthMax("login", utils.PLoginMaxLengthLimit))
var PasswordLoginMinLengthError = errors.New(lengthMin("login", utils.PLoginMinLengthLimit))
var PPasswordMaxLengthError = errors.New(lengthMax("password", utils.PPasswordMaxLengthLimit))
var PPasswordMinLengthError = errors.New(lengthMin("password", utils.PPasswordMinLengthLimit))
var EmailMaxLengthError = errors.New(lengthMax("email", utils.PEmailMinLengthLimit))
var EmailMinLengthError = errors.New(lengthMin("email", utils.PEmailMaxLengthLimit))