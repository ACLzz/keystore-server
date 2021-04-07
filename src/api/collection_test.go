package api

import (
	"fmt"
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/ACLzz/keystore-server/src/tests"
	"github.com/ACLzz/keystore-server/src/utils"
	"net/http"
	"testing"
)

func TestValidTitle(_t *testing.T) {
	testUserId := 7
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, _t)
	path := "collection/"
	url := fmt.Sprint(tests.BaseUrl, path)
	post := func(data map[string]interface{}, t *testing.T) ([]byte, *http.Response) {
		return tests.Post(url, data, t)
	}

	_t.Run("empty title", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.NoTitleError.Error())
		data := map[string]interface{}{"token": token}
		body, resp := post(data, t)

		tests.CheckResp(resp, body, 400, rightBody, t)
	})

	_t.Run("minimum title", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.CollectionTitleMinLengthError.Error())
		data := map[string]interface{}{"title": tests.BuildString(utils.CollectionTitleMinLengthLimit - 1), "token": token}
		body, resp := post(data, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("maximum title", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.CollectionTitleMaxLengthError.Error())
		data := map[string]interface{}{"title": tests.BuildString(utils.CollectionTitleMaxLengthLimit + 1), "token": token}
		body, resp := post(data, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("non-ascii title", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.CollectionLocaleError.Error())
		data := map[string]interface{}{"title": tests.BuildString(utils.PasswordMinLengthLimit+ 1) + "тест", "token": token}
		body, resp := post(data, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})
	
	tests.DeleteUser(testUserId, _t)
}


func TestCollectionAlreadyExists(t *testing.T) {
	testUserId := 8
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, t)
	user := *tests.GetUser(testUserId)
	testCollectionId := 1
	tests.CreateCollection(testCollectionId, user)
	
	path := "collection/"
	rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.CollectionExist.Error())
	url := fmt.Sprint(tests.BaseUrl, path)

	body, resp := tests.Post(url, map[string]interface{}{"title": tests.BuildTitle(testCollectionId), "token": token}, t)
	tests.CheckResp(resp, body, 422, rightBody, t)

	tests.DeleteCollection(testCollectionId, user, t)
	tests.DeleteUser(testUserId, t)
}

func TestCreateCollection(t *testing.T) {
	testUserId := 9
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, t)
	user := *tests.GetUser(testUserId)
	testCollectionId := 2
	path := "collection/"
	url := fmt.Sprint(tests.BaseUrl, path)
	title := tests.BuildTitle(testCollectionId)
	rightBody := ""

	body, resp := tests.Post(url, map[string]interface{}{"title": title, "token": token}, t)
	tests.CheckResp(resp, body, 201, rightBody, t)

	c := database.Collection{
		Title: title,
		User: user,
	}
	if !c.IsExist() {
		t.Error("Collection haven't created")
	} else {
		tests.DeleteCollection(testCollectionId, user, t)
	}

	tests.DeleteUser(testUserId, t)
}