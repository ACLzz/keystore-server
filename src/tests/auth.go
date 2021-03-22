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
	baseUser.Username = fmt.Sprintf("%s--%d", baseUser.Username, id)
	baseUser.Register()
}

func DeleteUser(id int, t *testing.T) {
	baseUser := BaseUser
	baseUser.Username = fmt.Sprintf("%s--%d", baseUser.Username, id)

	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()

	if exists := baseUser.IsExist(baseUser.Username); !exists {
		return
	}

	if tx := conn.Unscoped().Where("username = ?", baseUser.Username).Delete(database.User{}); tx.Error != nil {
		t.Errorf("Delete user: %v", tx.Error)
	}
}
