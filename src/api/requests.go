package api

import (
	"encoding/json"
	"fmt"
	"github.com/ACLzz/keystore-server/src/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func SendError(w http.ResponseWriter, err error, errCode int) {
	msg := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
	http.Error(w, msg, errCode)
}

func SendResp(w  http.ResponseWriter, resp map[string]interface{}, statusCode int) {
	jresp, err := json.Marshal(resp)
	if err != nil {
		log.Error("Error in json marshalling: ", err)
		SendError(w, errors.InternalError, 500)
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(jresp); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error in send response.")
	}
}


func ConvBody(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error in getting body, Register: %v\n", err)
		SendError(w, errors.InvalidBodyError, 400)
		return nil
	}

	if len(string(body[:])) <= 2 {
		SendError(w, errors.EmptyBodyError, 400)
		return nil
	}

	jbody := make(map[string]interface{})
	err = json.Unmarshal(body, &jbody)
	if err != nil {
		log.Error("Error in body unmarshalling: ", err)
		SendError(w, errors.InvalidBodyError, 400)
		return nil
	}
	return jbody
}
