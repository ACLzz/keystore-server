package api

import (
	"fmt"
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/ACLzz/keystore-server/src/tests"
	"gorm.io/gorm"
	"net/http"
	"testing"
	"time"
)

func TestEmptyBody(t *testing.T) {
	rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.EmptyBodyError.Error())
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)
	body, resp := tests.Post(url, map[string]interface{}{}, t)

	tests.CheckResp(resp, body, 400, rightBody, t)
}

func TestValidPassword(_t *testing.T) {
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)
	post := func(data map[string]interface{}, t *testing.T) ([]byte, *http.Response) {
		return tests.Post(url, data, t)
	}
	
	_t.Run("empty password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.NoPasswordError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit)}
		body, resp := post(data, t)
		
		tests.CheckResp(resp, body, 400, rightBody, t)
	})

	_t.Run("minimum password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PasswordMinLengthError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit),
			"password": tests.BuildString(errors.PasswordMinLengthLimit - 1)}
		body, resp := post(data, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})
	
	_t.Run("maximum password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PasswordMaxLengthError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit),
			"password": tests.BuildString(errors.PasswordMaxLengthLimit + 1)}
		body, resp := post(data, t)
		
		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("non-ascii password", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.PasswordLocaleError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit), "password":
		tests.BuildString(errors.PasswordMinLengthLimit + 1) + "тест"}
		body, resp := post(data, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})
}

func TestValidLogin(_t *testing.T) {
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)
	post := func(data map[string]interface{}, t *testing.T) ([]byte, *http.Response) {
		return tests.Post(url, data, t)
	}

	_t.Run("empty login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.NoLoginError.Error())
		data := map[string]interface{}{"password": tests.BuildString(errors.PasswordMinLengthLimit)}
		body, resp := post(data, t)

		tests.CheckResp(resp, body, 400, rightBody, t)
	})

	_t.Run("minimum login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.LoginMinLengthError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit - 1),
			"password": tests.BuildString(errors.PasswordMinLengthLimit)}
		body, resp := post(data, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("maximum login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.LoginMaxLengthError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMaxLengthLimit + 1),
			"password": tests.BuildString(errors.PasswordMinLengthLimit)}
		body, resp := post(data, t)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("non-ascii login", func(t *testing.T) {
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.LoginLocaleError.Error())
		data := map[string]interface{}{"login": tests.BuildString(errors.LoginMinLengthLimit) + "тест",
			"password": tests.BuildString(errors.PasswordMinLengthLimit)}
		body, resp := post(data, t)

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
		map[string]interface{}{"login": tests.BuildUsername(testUserId), "password": tests.BuildString(errors.PasswordMinLengthLimit)}, t)
	tests.CheckResp(resp, body, 422, rightBody, t)
	
	tests.DeleteUser(testUserId, t)
}

func TestRegister(t *testing.T) {
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)
	testUserId := 2
	username := tests.BuildUsername(testUserId)
	rightBody := ""

	body, resp := tests.Post(url,
		map[string]interface{}{"login": username, "password": tests.BuildString(errors.PasswordMinLengthLimit)}, t)
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
		body, resp := tests.Post(url, map[string]interface{}{"login": username, "password": tests.BaseUser.Password}, t)
		conn := database.GetConn()
		DB, _ := conn.DB()
		defer DB.Close()

		var user database.User
		var token database.Token
		conn.Where("username = ?", username).First(&user)
		conn.First(&token).Where("user_refer = ?", user.Id)
		rightBody := fmt.Sprintf("{\"token\":\"%s\"}", token.Token)

		tests.CheckResp(resp, body, 202, rightBody, t)
		tests.DeleteUser(testUserId, t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		body, resp := tests.Post(url, map[string]interface{}{"login": "this_user_not_exist", "password": tests.BaseUser.Password}, t)
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.InvalidCredentials.Error())

		tests.CheckResp(resp, body, 401, rightBody, t)
	})
}

func TestVerifyAuth(_t *testing.T) {
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)
	
	d := func(data map[string]interface{}, t *testing.T) ([]byte, *http.Response) {
		return tests.Delete(url, data, t)
	}
	
	_t.Run("no token", func(t *testing.T) {
		body, resp := d(map[string]interface{}{"token": ""}, t)
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.NoToken)
		
		tests.CheckResp(resp, body, 400, rightBody, t)
	})

	_t.Run("invalid token", func(t *testing.T) {
		body, resp := d(map[string]interface{}{"token": "invalidToken"}, t)
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.InvalidToken)

		tests.CheckResp(resp, body, 422, rightBody, t)
	})

	_t.Run("expired token", func(t *testing.T) {
		testUserId := 4
		_token := database.Token{
			Token:        "expiredToken",
			UserRefer:    "1",
			User:         database.User{Id: 1, Username: tests.BuildUsername(testUserId)},
			CreationDate: time.Now(),
			ExpireDate:   time.Now().Add(-24*time.Hour),
		}
		conn := database.GetConn()
		DB, _ := conn.DB()
		defer DB.Close()
		if tx := conn.Create(&_token); tx.Error != nil {
			t.Error(tx.Error)
		}

		body, resp := d(map[string]interface{}{"token": _token.Token}, t)
		rightBody := fmt.Sprintf("{\"error\":\"%s\"}\n", errors.ExpiredToken)

		tests.CheckResp(resp, body, 403, rightBody, t)
		conn.Delete(&_token)
		tests.DeleteUser(testUserId, t)
	})
}

func TestDeleteUser(t *testing.T) {
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)

	d := func(data map[string]interface{}) ([]byte, *http.Response) {
		return tests.Delete(url, data, t)
	}

	testUserId := 5
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, t)

	body, resp := d(map[string]interface{}{"token": token})
	rightBody := ""

	tests.CheckResp(resp, body, 200, rightBody, t)

	var user database.User
	username := tests.BuildUsername(testUserId)
	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	if tx := conn.First(&user).Where("username = ?", username); tx.Error != gorm.ErrRecordNotFound {
		t.Error("User haven't deleted")
	}
}

func TestUpdateUser(t *testing.T) {
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)

	p := func(data map[string]interface{}) ([]byte, *http.Response) {
		return tests.Put(url, data, t)
	}

	testUserId := 6
	updatedUsername := fmt.Sprint("updatedUsername", testUserId)
	updatedPassword := "updatedPassword"
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, t)

	body, resp := p(map[string]interface{}{"token": token, "login": updatedUsername, "password": updatedPassword})
	rightBody := ""

	tests.CheckResp(resp, body, 200, rightBody, t)

	var user database.User
	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	if tx := conn.First(&user).Where("username = ?", updatedUsername); tx.Error != nil {
		t.Error("User haven't updated username")
	}

	emptyUser := database.User{Password: updatedPassword}
	if user.Password != emptyUser.HashPassword() {
		fmt.Println(user.Password, emptyUser.Password)
		t.Error("Password haven't updated")
	}

	conn.Unscoped().Where("user_refer = ?", user.Username).Delete(database.Token{})

	if tx := conn.Unscoped().Where("username = ?", user.Username).Delete(database.User{}); tx.Error != nil {
		t.Errorf("Delete user: %v", tx.Error)
	}
}

func TestReadUser(t *testing.T) {
	path := "auth/"
	url := fmt.Sprint(tests.BaseUrl, path)

	testUserId := 7
	tests.RegisterUser(testUserId)
	token := tests.GetToken(testUserId, t)
	body, resp := tests.Get(url, map[string]interface{}{"token": token}, t)

	var user database.User
	conn := database.GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	if tx := conn.First(&user).Where("username = ?", tests.BuildUsername(testUserId)); tx.Error != nil {
		t.Error("User haven't updated username")
	}

	rightBody := fmt.Sprintf("{\"registered\":\"%s\",\"username\":\"%s\"}", user.RegistrationDate.Format("2006-01-02T15:04:05Z07:00"), user.Username)

	tests.CheckResp(resp, body, 200, rightBody, t)
	tests.DeleteUser(testUserId, t)
}
