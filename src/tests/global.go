package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ACLzz/keystore-server/src/utils"
	"io/ioutil"
	"net/http"
	"testing"
)

var Client = &http.Client{}
var BaseUrl = fmt.Sprintf("http://%s:%d/", utils.Config.Addr, utils.Config.Port)

func CheckResp(resp *http.Response, body []byte, code int, rightBody string, t *testing.T) {
	if resp.StatusCode != code {
		t.Errorf("status must be %d, but it is %d. error: %s", code, resp.StatusCode, body)
	} else if string(body) != rightBody {
		t.Errorf("response must be %s, but it is %s.", rightBody, string(body))
	}
}

func BuildString(length int) string {
	str := ""
	for i := 0; i<length; i++ {
		str += "a"
	}
	return str
}

func Post(url string, data map[string]interface{}, t *testing.T) ([]byte, *http.Response) {
	jvalues, _ := json.Marshal(data)
	resp, err := Client.Post(url, "application/json", bytes.NewBuffer(jvalues))
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, resp
}

func Delete(url string, data map[string]interface{}, t *testing.T) ([]byte, *http.Response) {
	return customRequest("DELETE", url, data, t)
}

func Put(url string, data map[string]interface{}, t *testing.T) ([]byte, *http.Response) {
	return customRequest("PUT", url, data, t)
}

func customRequest(method string, url string, data map[string]interface{}, t *testing.T) ([]byte, *http.Response) {
	jvalues, _ := json.Marshal(data)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jvalues))
	if err != nil {
		t.Error(err)
	}

	resp, err := Client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, resp
}
