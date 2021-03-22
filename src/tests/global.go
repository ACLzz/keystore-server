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

func CheckResp(resp *http.Response, body []byte, code int, rightBody string) string {
	if resp.StatusCode != code {
		return fmt.Sprintf("status must be %d, but it is %d. error: %s", code, resp.StatusCode, body)
	} else if string(body) != rightBody {
		return fmt.Sprintf("response must be %s, but it is %s.", rightBody, string(body))
	}

	return ""
}

func BuildString(length int) string {
	str := ""
	for i := 0; i<length; i++ {
		str += ""
	}
	return str
}

func Post(url string, data map[string]interface{}, subTestName string, t *testing.T) ([]byte, *http.Response) {
	jvalues, _ := json.Marshal(data)
	resp, err := Client.Post(url, "application/json", bytes.NewBuffer(jvalues))
	if err != nil {
		t.Errorf("Error in %s: %v", subTestName, err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, resp
}
