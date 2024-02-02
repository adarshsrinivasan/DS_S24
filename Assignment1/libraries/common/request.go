package common

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Credentials struct {
	Name     string `json:"name,omitempty" bson:"name" bun:"name,notnull"`
	UserName string `json:"userName,omitempty" bson:"userName" bun:"userName,notnull,unique"`
	Password string `json:"password,omitempty" bson:"password" bun:"password,notnull,unique"`
}

type ClientRequest struct {
	SessionID string
	Service   string
	UserType  UserType
	Body      []byte
}

func (req *ClientRequest) SerializeRequest() ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(req)
	if err != nil {
		log.Fatal("Encode error:", err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (req *ClientRequest) DeserializeRequest(data []byte) error {
	var buffer bytes.Buffer
	buffer.Write(data)
	decoder := gob.NewDecoder(&buffer)
	if err := decoder.Decode(&req); err != nil {
		log.Fatal("Decode error:", err)
		return err
	}
	return nil
}
