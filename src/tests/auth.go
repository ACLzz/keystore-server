package tests

import (
	"fmt"
	"github.com/ACLzz/keystore-server/src/database"
	"testing"
)

var BaseUser = database.User{
	Username:         "tester",
	Password:         "testing_password",
}

func RegisterUser(id int) {
	baseUser := BaseUser
	baseUser.Username = BuildUsername(id)
	baseUser.Register()
}

func DeleteUser(id int, t *testing.T) {
	baseUser := BaseUser
	baseUser.Username = BuildUsername(id)

	conn := database.GetConn()

	if !baseUser.IsExist() {
		t.Errorf("User with id %d doesn't exist to delete him.", id)
		return
	}

	conn.Unscoped().Where("user_refer = ?", baseUser.Username).Delete(database.Token{})

	if tx := conn.Unscoped().Where("username = ?", baseUser.Username).Delete(database.User{}); tx.Error != nil {
		t.Errorf("Delete user: %v", tx.Error)
	}
}

func BuildUsername(id int) string {
	return fmt.Sprint(BaseUser.Username, id)
}
