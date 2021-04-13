package tests

import (
	"fmt"
	"github.com/ACLzz/keystore-server/src/database"
	"strconv"
	"testing"
)

var BasePassword = database.Password{
	Email:           "test-email@mail.com",
	Login:           "testlogin",
	Password:        "testpassword",
}

func CreatePassword(id int, collection *database.Collection) {
	basePassword := BasePassword
	basePassword.Title = BuildTitle(id)
	basePassword.Collection = *collection
	basePassword.CollectionRefer = fmt.Sprint(collection.Id)
	userId, _ := strconv.Atoi(collection.UserRefer)
	basePassword.Add(collection.Title, userId)
}

func BuildPasswordEntityString(collection *database.Collection, testPasswordId int, t *testing.T) string {
	passwords := collection.FetchPasswords()
	passwordTitle := BuildTitle(testPasswordId)
	for _, password := range passwords {
		if password.Title == passwordTitle {
			return fmt.Sprintf("{\"id\":%d,\"title\":\"%s\"}", password.Id, password.Title)
		}
	}
	t.Error("No password with title ", passwordTitle, " was found")
	return ""
}
