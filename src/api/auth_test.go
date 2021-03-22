package api

import (
	"fmt"
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/ACLzz/keystore-server/src/tests"
	"net/http"
	"testing"
)

func TestEmptyBody(t *testing.T) {
	rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.EmptyBodyError.Error())
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)
	body, resp := tests.Post(url, map[string]interface{}{}, "EmptyBody", t)

	tests.CheckResp(resp, body, 400, rightBody, t)
}

func TestValidPassword(t *testing.T) {
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)
	post := func(data map[string]interface{}, subTestName string) ([]byte, *http.Response) {
		return tests.Post(url, data, fmt.Sprint("TestValidPassword.", subTestName), t)
	}
	
	t.Run("empty password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.NoPasswordError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit)}
		body, resp := post(data, "empty_password")
		
		tests.CheckResp(resp, body, 400, rightBody, t)
	})

	t.Run("minimum password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PasswordMinLengthError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit),
			"password": tests.BuildString(errors.PasswordMinLengthLimit - 1)}
		body, resp := post(data, "min_password")

		tests.CheckResp(resp, body, 422, rightBody, t)
	})
	
	t.Run("maximum password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PasswordMaxLengthError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit),
			"password": tests.BuildString(errors.PasswordMaxLengthLimit + 1)}
		body, resp := post(data, "max_password")
		
		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	t.Run("non-ascii password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PasswordLocaleError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit), "password":
		tests.BuildString(errors.PasswordMinLengthLimit + 1) + "тест"}
		body, resp := post(data, "non-ascii_password")

		tests.CheckResp(resp, body, 422, rightBody, t)
	})
}

func TestValidLogin(t *testing.T) {
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)
	post := func(data map[string]interface{}, subTestName string) ([]byte, *http.Response) {
		return tests.Post(url, data, fmt.Sprint("TestValidLogin.", subTestName), t)
	}

	t.Run("empty login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.NoLoginError.Error())
		data := map[string]interface{}{"password": tests.BuildString(errors.PasswordMinLengthLimit)}
		body, resp := post(data, "empty_login")

		tests.CheckResp(resp, body, 400, rightBody, t)
	})

	t.Run("minimum login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.LoginMinLengthError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit - 1),
			"password": tests.BuildString(errors.PasswordMinLengthLimit)}
		body, resp := post(data, "min_login")

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	t.Run("maximum login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.LoginMaxLengthError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMaxLengthLimit + 1),
			"password": tests.BuildString(errors.PasswordMinLengthLimit)}
		body, resp := post(data, "max_login")

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	t.Run("non-ascii login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.LoginLocaleError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit) + "тест",
			"password": tests.BuildString(errors.PasswordMinLengthLimit)}
		body, resp := post(data, "non-ascii_login")

		tests.CheckResp(resp, body, 422, rightBody, t)
	})
}

func TestUserAlreadyExists(t *testing.T) {
	testUserId := 1
	tests.RegisterUser(testUserId)
	path := "auth/"
	rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.UserExists.Error())
	url := fmt.Sprint(tests.BaseUrl, path)

	body, resp := tests.Post(url,
		map[string]interface{}{"login": tests.BuildUsername(testUserId), "password": tests.BuildString(errors.PasswordMinLengthLimit)},
		"UserAlreadyExists", t)
	tests.CheckResp(resp, body, 422, rightBody, t)
	
	tests.DeleteUser(testUserId, t)
}

func TestRegister(t *testing.T) {
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)
	testUserId := 2
	username := tests.BuildUsername(testUserId)
	rightBody := "{\"ok\":true}"

	body, resp := tests.Post(url,
		map[string]interface{}{"login": username, "password": tests.BuildString(errors.PasswordMinLengthLimit)},
		"Register", t)
	tests.CheckResp(resp, body, 201, rightBody, t)

	u := database.User{
		Username: username,
	}
	if !u.IsExist() {
		t.Error("User haven't registered")
	} else {
		tests.DeleteUser(testUserId, t)
	}
}

func TestLogin(t *testing.T) {
	path := "auth/login"
	url := fmt.Sprint(tests.BaseUrl, path)

	t.Run("get token", func(t *testing.T) {
		testUserId := 3
		username := tests.BuildUsername(testUserId)
		tests.RegisterUser(testUserId)
		body, resp := tests.Post(url, map[string]interface{}{"login": username, "password": tests.BaseUser.Password}, "TestLogin", t)
		conn := database.GetConn()

		var user database.User
		var token database.Token
		conn.Where("username = ?", username).First(&user)
		conn.First(&token).Where("user_refer = ?", user.Id)
		rightBody := fmt.Sprintf("{\"token\":\"%s\"}", token.Token)

		tests.CheckResp(resp, body, 202, rightBody, t)
		tests.DeleteUser(testUserId, t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		body, resp := tests.Post(url, map[string]interface{}{"login": "this_user_not_exist", "password": tests.BaseUser.Password}, "TestLogin", t)

		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.InvalidCredentials.Error())

		tests.CheckResp(resp, body, 401, rightBody, t)
	})
}
