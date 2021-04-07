package errors

import (
	"errors"
	"github.com/ACLzz/keystore-server/src/utils"
)

var CollectionTitleMaxLengthError = errors.New(lengthMax("title", utils.CollectionTitleMaxLengthLimit))
var CollectionTitleMinLengthError = errors.New(lengthMin("title", utils.CollectionTitleMinLengthLimit))

var CollectionLocaleError = errors.New(localeError("title"))

var NoTitleError = errors.New(noField("title"))

var CollectionExist = errors.New("collection with that title already exist")
var CollectionNotExist = errors.New("collection with that title does not exist")