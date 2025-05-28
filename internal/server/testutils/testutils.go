package testutils

import (
	"encoding/json"
	"go-keeper/internal/common/protocol"
)

func MustCreateCredsJSON(username, password string) string {
	creds := protocol.Creds{
		Username: username,
		Password: password,
	}
	jsonCreds, err := json.Marshal(creds)
	if err != nil {
		panic(err)
	}
	return string(jsonCreds)
}
