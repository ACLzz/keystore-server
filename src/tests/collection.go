package tests

import (
	"fmt"
	"github.com/ACLzz/keystore-server/src/database"
	"testing"
)

var BaseCollection = database.Collection{
	Title: "title",
}

func CreateCollection(id int, user database.User) {
	baseCollection := BaseCollection
	baseCollection.Title = BuildTitle(id)
	baseCollection.User = user
	baseCollection.Add()
}

func DeleteCollection(id int, user database.User, t *testing.T) {
	baseCollection := BaseCollection
	baseCollection.Title = BuildTitle(id)
	baseCollection.User = user

	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()

	if !baseCollection.IsExist() {
		t.Errorf("Collection with id %d for user %s doesn't exist to delete it.", id, baseCollection.User.Username)
		return
	}

	if tx := conn.Unscoped().Where("title = ? and user_refer = ?", baseCollection.Title, baseCollection.User.Id).Delete(database.Collection{});
	tx.Error != nil {
		t.Errorf("Delete collection: %v", tx.Error)
	}
}

func BuildTitle(id int) string {
	return fmt.Sprint(BaseUser.Username, id)
}

