package api

import (
	"fmt"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/ACLzz/keystore-server/src/tests"
	"testing"
)

func TestPasswordMiddleware(t *testing.T) {
	testUserId := 13
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, t)
	path := "collection/notExist/"
	url := fmt.Sprint(tests.BaseUrl, path)
	rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.CollectionNotExist.Error())
	body, resp := tests.Post(url,
		map[string]interface{}{"title": "whatever", "login": "whatever", "password": "whatever", "token": token}, t)

	tests.CheckResp(resp, body, 404, rightBody, t)
	tests.DeleteUser(testUserId, t)
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
			map[string]interface{}{"login": login, "password": password, "token": token}, t)

		tests.CheckResp(resp, body, 400, rightBody, t)
	})

	_t.Run("no password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.NoPasswordError.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"title": title, "login": login, "token": token}, t)

		tests.CheckResp(resp, body, 400, rightBody, t)
	})
	_t.Run("no login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.NoLoginError.Error())
		body, resp := tests.Post(url,
			map[string]interface{}{"title": title, "password": password, "token": token}, t)

		tests.CheckResp(resp, body, 400, rightBody, t)
	})
	_t.Run("create password", func(t *testing.T) {
		rightBody := ""
		body, resp := tests.Post(url,
			map[string]interface{}{"title": title,"login": login, "password": password, "token": token}, t)

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
			map[string]interface{}{"token": token}, t)

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
			map[string]interface{}{"token": token}, t)

		tests.CheckResp(resp, body, 200, rightBody, t)
	})

	tests.DeleteCollection(testCollectionId, *user, _t)
	tests.DeleteUser(testUserId, _t)
}