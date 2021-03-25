package api

import (
	"github.com/ACLzz/keystore-server/src/utils"
	"net/http"
)

func ShutdownServer(w http.ResponseWriter, _ *http.Request) {
	utils.EndCh <- 1
	SendResp(w, nil, 200)
}
