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
	post := func(data map[string]interface{}, headers map[string]string, t *testing.T) ([]byte, *http.Response) {
		return tests.Post(url, data, headers, t)
	}

	_t.Run("empty title", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.NoTitleError.Error())
		body, resp := post(map[string]interface{}{"noNeedKey": "yep"}, map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 400, rightBody, t)
	})

	_t.Run("minimum title", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.CollectionTitleMinLengthError.Error())
		data := map[string]interface{}{"title": tests.BuildString(utils.CollectionTitleMinLengthLimit - 1)}
		body, resp := post(data, map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("maximum title", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.CollectionTitleMaxLengthError.Error())
		data := map[string]interface{}{"title": tests.BuildString(utils.CollectionTitleMaxLengthLimit + 1)}
		body, resp := post(data, map[string]string{"Authorization": token}, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("non-ascii title", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.CollectionLocaleError.Error())
		data := map[string]interface{}{"title": tests.BuildString(utils.PasswordMinLengthLimit+ 1) + "ัะตัั"}
		body, resp := post(data, map[string]string{"Authorization": token}, t)

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

	body, resp := tests.Post(url, map[string]interface{}{"title": tests.BuildTitle(testCollectionId)}, map[string]string{"Authorization": token}, t)
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

	body, resp := tests.Post(url, map[string]interface{}{"title": title}, map[string]string{"Authorization":  token}, t)
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

func TestFetchCollections(t *testing.T) {
	testUserId := 10
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, t)
	user := *tests.GetUser(testUserId)
	testCollection1Id := 3
	testCollection2Id := 4
	tests.CreateCollection(testCollection1Id, user)
	tests.CreateCollection(testCollection2Id, user)

	path := "collection/"
	rightBody := fmt.Sprintf("[\"%s\",\"%s\"]", tests.BuildTitle(testCollection1Id), tests.BuildTitle(testCollection2Id))
	url := fmt.Sprint(tests.BaseUrl, path)

	body, resp := tests.Get(url, map[string]string{"Authorization": token}, t)
	tests.CheckResp(resp, body, 200, rightBody, t)

	tests.DeleteCollection(testCollection1Id, user, t)
	tests.DeleteCollection(testCollection2Id, user, t)
	tests.DeleteUser(testUserId, t)
}

func TestUpdateCollection(_t *testing.T) {
	testUserId := 11
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, _t)
	user := *tests.GetUser(testUserId)
	_path := "collection/"
	
	_t.Run("collection does not exist", func(t *testing.T) {
		path := fmt.Sprint(_path, "not_exist")
		url := fmt.Sprint(tests.BaseUrl, path)
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.CollectionNotExist.Error())

		body, resp := tests.Put(url, map[string]interface{}{"title": "test"}, map[string]string{"Authorization": token}, t)
		tests.CheckResp(resp, body, 404, rightBody, t)
	})
	
	_t.Run("collection already exist", func(t *testing.T) {
		testCollection1Id := 5
		testCollection2Id := 6
		tests.CreateCollection(testCollection1Id, user)
		tests.CreateCollection(testCollection2Id, user)
		title := tests.BuildTitle(testCollection1Id)
		path := fmt.Sprint(_path, title)
		url := fmt.Sprint(tests.BaseUrl, path)
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.CollectionExist.Error())

		body, resp := tests.Put(url, map[string]interface{}{"title": tests.BuildTitle(testCollection2Id)}, map[string]string{"Authorization": token}, t)
		tests.CheckResp(resp, body, 422, rightBody, t)

		tests.DeleteCollection(testCollection1Id, user, t)
		tests.DeleteCollection(testCollection2Id, user, t)
	})
	
	_t.Run("update", func(t *testing.T) {
		testCollectionId := 7
		tests.CreateCollection(testCollectionId, user)
		title := tests.BuildTitle(testCollectionId)
		newTitle := fmt.Sprint(title, "updated")
		path := fmt.Sprint(_path, title)
		url := fmt.Sprint(tests.BaseUrl, path)
		rightBody := ""

		body, resp := tests.Put(url, map[string]interface{}{"title": newTitle}, map[string]string{"Authorization": token}, t)
		tests.CheckResp(resp, body, 200, rightBody, t)

		conn := database.GetConn()
		DB, _ := conn.DB()
		defer DB.Close()
		
		if tx := conn.Unscoped().Where("title = ? and user_refer = ?", newTitle, user.Id).Delete(database.Collection{}); tx.Error != nil {
			t.Error(tx.Error)
		}
	})

	tests.DeleteUser(testUserId, _t)
}

func TestDeleteCollection(_t *testing.T) {
	_path := "collection/"
	testUserId := 12
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, _t)
	user := *tests.GetUser(testUserId)

	_t.Run("collection does not exist", func(t *testing.T) {
		path := fmt.Sprint(_path, "not_exist")
		url := fmt.Sprint(tests.BaseUrl, path)
		body, resp := tests.Delete(url, map[string]string{"Authorization": token}, t)
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.CollectionNotExist.Error())

		tests.CheckResp(resp, body, 404, rightBody, t)	
	})
	
	_t.Run("delete collection", func(t *testing.T) {
		testCollectionId := 7
		tests.CreateCollection(testCollectionId, user)
		title := tests.BuildTitle(testCollectionId)
		path := fmt.Sprint(_path, title)
		url := fmt.Sprint(tests.BaseUrl, path)
		body, resp := tests.Delete(url, map[string]string{"Authorization": token}, t)
		rightBody := ""

		tests.CheckResp(resp, body, 200, rightBody, t)
		
		c := database.Collection{Title: title, User: user}
		
		if c.IsExist() {
			tests.DeleteCollection(testCollectionId, user, t)
			t.Error("Collection haven't deleted")
		}
	})
	
	tests.DeleteUser(testUserId, _t)
}

func TestListCollection(_t *testing.T) {
	testUserId := 15
	tests.RegisterUser(testUserId)
	user := *tests.GetUser(testUserId)
	token := tests.GetToken(testUserId, _t)
	_path := fmt.Sprint(tests.BaseUrl, "collection/")
	
	_t.Run("empty collection list", func(t *testing.T) {
		testCollectionId := 8
		tests.CreateCollection(testCollectionId, user)
		url := fmt.Sprint(_path, tests.BuildTitle(testCollectionId))
		body, resp := tests.Get(url, map[string]string{"Authorization": token}, t)
		rightBody := "[]"

		tests.CheckResp(resp, body, 200, rightBody, t)
		tests.DeleteCollection(testCollectionId, user, t)
	})
	
	_t.Run("list collection", func(t *testing.T) {
		testCollectionId := 9
		tests.CreateCollection(testCollectionId, user)
		title := tests.BuildTitle(testCollectionId)
		url := fmt.Sprint(_path, title)
		collection := GetCollection(title, &user)
		tests.CreatePassword(1, *collection)
		tests.CreatePassword(2, *collection)
		tests.CreatePassword(3, *collection)

		body, resp := tests.Get(url, map[string]string{"Authorization": token}, t)
		rightBody := fmt.Sprintf("[%s,%s,%s]",
			tests.BuildPasswordEntityString(collection, 1, t),
			tests.BuildPasswordEntityString(collection, 2, t),
			tests.BuildPasswordEntityString(collection, 3, t))

		tests.CheckResp(resp, body, 200, rightBody, t)
		tests.DeleteCollection(testCollectionId, user, t)
	})

	tests.DeleteUser(testUserId, _t)
}