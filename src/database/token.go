package database

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/ACLzz/keystore-server/src/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

func (t *Token) genToken() {
	conn := GetConn()
	defer conn.Commit()
	
	// Prepare token dates
	now := time.Now()
	loc, _ := time.LoadLocation(utils.Config.Timezone)
	t.CreationDate = now.In(loc)
	t.ExpireDate = now.Add(time.Duration(utils.Config.TokenLifetime * 3600000000000)).In(loc)

	jstruct, err := json.Marshal(t)
	if err != nil {
		log.Error("Error in token generator: ", err.Error())
	}
	t.Token = fmt.Sprintf("%x", sha256.Sum256(jstruct))

	if tx := conn.Create(t); tx.Error != nil {
		log.Error("Error in adding token to database: ", tx.Error.Error())
	}
	conn.Commit()
}
