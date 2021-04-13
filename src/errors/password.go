package errors

import "errors"

var PasswordNotExist = errors.New("password with that id does not exist")
var InvalidPasswordId = errors.New("you sent invalid password id")