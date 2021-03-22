package api

import (
	"fmt"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/ACLzz/keystore-server/src/tests"
	"net/http"
	"testing"
)

func TestEmptyBody(t *testing.T) {
	rightBody := fmt.Sprintf("{\"error\": \"%s\"}\n", errors.EmptyBodyError.Error())
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)
	body, resp := tests.Post(url, map[string]interface{}{}, "EmptyBody", t)

	tests.CheckResp(resp, body, 400, rightBody)
}

func TestRegisterPassword(t *testing.T) {
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)
	post := func(data map[string]interface{}, subTestName string) ([]byte, *http.Response) {
		return tests.Post(url, data, fmt.Sprint("RegisterTestPassword.", subTestName), t)
	}
	
	t.Run("empty password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\": \"%s\"}\n", errors.NoPasswordError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit)}
		body, resp := post(data, "empty_password")
		
		tests.CheckResp(resp, body, 400, rightBody)
	})

	t.Run("minimum password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\": \"%s\"}\n", errors.PasswordMinLengthError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit),
			"password": tests.BuildString(errors.PasswordMinLengthLimit - 1)}
		body, resp := post(data, "min_password")

		tests.CheckResp(resp, body, 422, rightBody)
	})
	
	t.Run("maximum password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\": \"%s\"}\n", errors.PasswordMaxLengthError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit),
			"password": tests.BuildString(errors.PasswordMaxLengthLimit + 1)}
		body, resp := post(data, "max_password")
		
		tests.CheckResp(resp, body, 422, rightBody)
	})

	t.Run("non-ascii password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\": \"%s\"}\n", errors.PasswordLocaleError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit), "password":
		tests.BuildString(errors.PasswordMinLengthLimit + 1) + "тест"}
		body, resp := post(data, "non-ascii_password")

		tests.CheckResp(resp, body, 422, rightBody)
	})
}

func TestRegisterLogin(t *testing.T) {
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)
	post := func(data map[string]interface{}, subTestName string) ([]byte, *http.Response) {
		return tests.Post(url, data, fmt.Sprint("RegisterTestLogin.", subTestName), t)
	}

	t.Run("empty login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\": \"%s\"}\n", errors.NoLoginError.Error())
		data := map[string]interface{}{"password": tests.BuildString(errors.PasswordMinLengthLimit)}
		body, resp := post(data, "empty_login")

		tests.CheckResp(resp, body, 400, rightBody)
	})

	t.Run("minimum login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\": \"%s\"}\n", errors.LoginMinLengthError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit - 1),
			"password": tests.BuildString(errors.PasswordMinLengthLimit)}
		body, resp := post(data, "min_login")

		tests.CheckResp(resp, body, 422, rightBody)
	})

	t.Run("maximum login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\": \"%s\"}\n", errors.LoginMaxLengthError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMaxLengthLimit + 1),
			"password": tests.BuildString(errors.PasswordMinLengthLimit)}
		body, resp := post(data, "max_login")

		tests.CheckResp(resp, body, 422, rightBody)
	})

	t.Run("non-ascii login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\": \"%s\"}\n", errors.LoginLocaleError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit) + "тест",
			"password": tests.BuildString(errors.PasswordMinLengthLimit)}
		body, resp := post(data, "non-ascii_login")

		tests.CheckResp(resp, body, 422, rightBody)
	})
}

func TestUserAlreadyExists(t *testing.T) {
	testUserId := 1
	tests.RegisterUser(testUserId)
	path := "auth/"
	rightBody := fmt.Sprintf("{\"error\": \"%s\"}\n", errors.UserExists.Error())
	url := fmt.Sprint(tests.BaseUrl, path)

	body, resp := tests.Post(url,
		map[string]interface{}{"login": tests.BaseUser.Username + "--1", "password": tests.BuildString(errors.PasswordMinLengthLimit)},
		"UserAlreadyExists", t)
	tests.CheckResp(resp, body, 422, rightBody)
	
	tests.DeleteUser(testUserId, t)
}
