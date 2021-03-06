package api

import (
	"fmt"
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/ACLzz/keystore-server/src/tests"
	"github.com/ACLzz/keystore-server/src/utils"
	"gorm.io/gorm"
	"testing"
)

func TestPasswordMiddleware(_t *testing.T) {
	testUserId := 13
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, _t)
	_path := "collection/"
	
	_t.Run("not exists collection", func(t *testing.T) {
		path := fmt.Sprint(_path, "not_exist")
		url := fmt.Sprint(tests.BaseUrl, path)
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.CollectionNotExist.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"title": "whatever", "login": "whatever", "password": "whatever"}, map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 404, rightBody, t)
	})

	_t.Run("non-int password", func(t *testing.T) {
		testCollectionId := 9
		user := *GetUser(token)
		tests.CreateCollection(testCollectionId, user)
		path := fmt.Sprint(_path, tests.BuildTitle(testCollectionId),"/non-int")
		url := fmt.Sprint(tests.BaseUrl, path)

		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.InvalidPasswordId)
		body, resp := tests.Get(url,
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
		tests.DeleteCollection(testCollectionId, user, t)
	})

	tests.DeleteUser(testUserId, _t)
}

func TestValidPasswordFields(_t *testing.T) {
	testUserId := 17
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, _t)
	user := tests.GetUser(testUserId)
	testCollectionId := 12
	tests.CreateCollection(testCollectionId, *user)

	title := "test_title"
	login := "test_login"
	password := "test_password"
	path := fmt.Sprint("collection/", tests.BuildTitle(testCollectionId), "/")
	url := fmt.Sprint(tests.BaseUrl, path)

	_t.Run("minimum title", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PasswordTitleMinLengthError.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"title": tests.BuildString(utils.PLoginMinLengthLimit - 1),"login": login, "password": password},
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("maximum title", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PasswordTitleMaxLengthError.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"title": tests.BuildString(utils.PasswordTitleMaxLengthLimit + 1),"login": login, "password": password},
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("minimum login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PasswordLoginMinLengthError.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"title": title,"login": tests.BuildString(utils.PLoginMinLengthLimit - 1), "password": password},
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("maximum login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PasswordLoginMaxLengthError.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"title": title,"login": tests.BuildString(utils.PLoginMaxLengthLimit + 1), "password": password},
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("minimum password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PPasswordMinLengthError.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"title": title,"login": login, "password": tests.BuildString(utils.PPasswordMinLengthLimit - 1)},
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("maximum password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PPasswordMaxLengthError.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"title": title,"login": login, "password": tests.BuildString(utils.PPasswordMaxLengthLimit + 1)},
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("maximum email", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.EmailMaxLengthError.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"title": title,"login": login, "password": password, "email": tests.BuildString(utils.PEmailMaxLengthLimit + 1)},
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	tests.DeleteCollection(testCollectionId, *user, _t)
	tests.DeleteUser(testUserId, _t)
}

func TestCreatePassword(_t *testing.T) {
	testUserId := 14
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, _t)
	user := tests.GetUser(testUserId)
	testCollectionId := 7
	tests.CreateCollection(testCollectionId, *user)
	
	title := "test_title"
	login := "test_login"
	password := "test_password"
	path := fmt.Sprint("collection/", tests.BuildTitle(testCollectionId), "/")
	url := fmt.Sprint(tests.BaseUrl, path)
	
	_t.Run("no title", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.NoTitleError.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"login": login, "password": password}, map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 400, rightBody, t)
	})

	_t.Run("no password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.NoPasswordError.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"title": title, "login": login}, map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 400, rightBody, t)
	})
	_t.Run("no login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.NoLoginError.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"title": title, "password": password}, map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 400, rightBody, t)
	})
	_t.Run("create password", func(t *testing.T) {
		rightBody := ""
		body, resp := tests.Post(url,
			map[string]interface{}{"title": title,"login": login, "password": password},
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 201, rightBody, t)
	})
	
	tests.DeleteCollection(testCollectionId, *user, _t)
	tests.DeleteUser(testUserId, _t)
}

func TestReadPassword(_t *testing.T) {
	testUserId := 14
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, _t)
	user := tests.GetUser(testUserId)
	testCollectionId := 8
	tests.CreateCollection(testCollectionId, *user)
	collection := GetCollection(tests.BuildTitle(testCollectionId), user)
	_path := fmt.Sprint("collection/", tests.BuildTitle(testCollectionId), "/")
	
	_t.Run("password not exists", func(t *testing.T) {
		path := fmt.Sprint(_path, "0")
		url := fmt.Sprint(tests.BaseUrl, path)

		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PasswordNotExist)
		body, resp := tests.Get(url,
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 404, rightBody, t)
	})
	
	_t.Run("read password", func(t *testing.T) {
		tests.CreatePassword(1, *collection)
		password := tests.GetPassword(collection, 1, t)
		path := fmt.Sprint(_path, password.Id)
		url := fmt.Sprint(tests.BaseUrl, path)

		rightBody := fmt.Sprintf("{\"email\":\"%s\",\"login\":\"%s\",\"password\":\"%s\",\"title\":\"%s\"}",
			tests.BasePassword.Email, tests.BasePassword.Login, tests.BasePassword.Password, tests.BuildTitle(1))
		body, resp := tests.Get(url,
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 200, rightBody, t)
	})

	tests.DeleteCollection(testCollectionId, *user, _t)
	tests.DeleteUser(testUserId, _t)
}

func TestDeletePassword(_t *testing.T) {
	testUserId := 15
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, _t)
	user := tests.GetUser(testUserId)
	testCollectionId := 10
	tests.CreateCollection(testCollectionId, *user)
	collection := GetCollection(tests.BuildTitle(testCollectionId), user)
	_path := fmt.Sprint("collection/", tests.BuildTitle(testCollectionId), "/")

	_t.Run("password not exists", func(t *testing.T) {
		path := fmt.Sprint(_path, "0")
		url := fmt.Sprint(tests.BaseUrl, path)

		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PasswordNotExist)
		body, resp := tests.Delete(url,
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 404, rightBody, t)
	})

	_t.Run("delete password", func(t *testing.T) {
		tests.CreatePassword(1, *collection)
		password := tests.GetPassword(collection, 1, t)
		path := fmt.Sprint(_path, password.Id)
		url := fmt.Sprint(tests.BaseUrl, path)

		rightBody := ""
		body, resp := tests.Delete(url,
			map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 200, rightBody, t)
		
		conn := database.GetConn()
		DB, _ := conn.DB()
		defer DB.Close()

		var psswd database.Password
		if tx := conn.Where("title = ? AND collection_refer = ?", password.Title, password.CollectionRefer).First(&psswd);
		tx.Error != gorm.ErrRecordNotFound {
			t.Error("Password haven't deleted")
		}
	})

	tests.DeleteCollection(testCollectionId, *user, _t)
	tests.DeleteUser(testUserId, _t)
}

func TestUpdatePassword(_t *testing.T) {
	testUserId := 16
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, _t)
	user := tests.GetUser(testUserId)
	testCollectionId1 := 11
	testCollectionId2 := 30
	tests.CreateCollection(testCollectionId1, *user)
	tests.CreateCollection(testCollectionId2, *user)
	collection1 := GetCollection(tests.BuildTitle(testCollectionId1), user)
	collection2 := GetCollection(tests.BuildTitle(testCollectionId2), user)
	_path := fmt.Sprint("collection/", tests.BuildTitle(testCollectionId1), "/")

	_t.Run("update partially", func(t *testing.T) {
		tests.CreatePassword(1, *collection1)
		passwordId := tests.GetPassword(collection1, 1, t).Id
		path := fmt.Sprint(_path, passwordId)
		url := fmt.Sprint(tests.BaseUrl, path)
		rightBody := ""
		newTitle := "password's new title"
		newLogin := "new login"

		body, resp := tests.Put(url, map[string]interface{}{"title": newTitle, "login": newLogin},
		map[string]string{"Authorization": token}, t)
		tests.CheckResp(resp, body, 200, rightBody, t)

		conn := database.GetConn()
		DB, _ := conn.DB()
		defer DB.Close()
		
		var password database.Password
		if tx := conn.Unscoped().Where("title = ? and login = ? and collection_refer = ?", newTitle, newLogin, collection1.Id).
			First(&password); tx.Error != nil {
			t.Error(tx.Error)
		}
		
		if password.Password != tests.BasePassword.Password {
			t.Error("password row was overwritten with blank password field")
		}
	})

	_t.Run("update full", func(t *testing.T) {
		tests.CreatePassword(2, *collection1)
		passwordId := tests.GetPassword(collection1, 2, t).Id
		path := fmt.Sprint(_path, passwordId)
		url := fmt.Sprint(tests.BaseUrl, path)
		rightBody := ""
		
		newTitle := "password's new title2"
		newLogin := "new login2"
		newPassword := "newPassword"
		newEmail := "newEmail"

		body, resp := tests.Put(url, map[string]interface{}{
			"title": newTitle, "login": newLogin,"password": newPassword, "email": newEmail, "collection": tests.BuildTitle(testCollectionId2)},
			map[string]string{"Authorization": token}, t)
		tests.CheckResp(resp, body, 200, rightBody, t)

		conn := database.GetConn()
		DB, _ := conn.DB()
		defer DB.Close()

		var password database.Password
		if tx := conn.Unscoped().Where("title = ? and login = ? and password = ? and email = ? and collection_refer = ?",
			newTitle, newLogin, newPassword, newEmail, collection2.Id).
			First(&password); tx.Error != nil {
			t.Error(tx.Error)
		}
	})

	tests.DeleteCollection(testCollectionId1, *user, _t)
	tests.DeleteCollection(testCollectionId2, *user, _t)
	tests.DeleteUser(testUserId, _t)
}