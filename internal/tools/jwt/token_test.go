package jwt

import (
	"encoding/json"
	"testing"
	"time"
)

type UserPub struct {
	ID    *int   `json:"user_id"`
	Email string `json:"email"`
}

func TestEncode(t *testing.T) {
	secret := "this is screct"
	id := 1
	userPub := UserPub{
		ID:    &id,
		Email: "Murphy@jwt.com",
	}
	payload := Payload{
		Sub: "123",
		Exp: time.Now().Unix() + 100000,
		Pub: userPub,
	}
	token, err := Encode(payload, secret)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	payloadDecoded, err := Decode(token, secret)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if payload.Exp != payloadDecoded.Exp {
		t.Errorf("Error: payload data not matched")
	}
	userPubRes := UserPub{}
	userPubBytes, _ := json.Marshal(payloadDecoded.Pub)
	json.Unmarshal(userPubBytes, &userPubRes)
	if userPub.Email != userPubRes.Email || *userPub.ID != *userPubRes.ID {
		t.Errorf("Error: payload pub not matched")
	}
}
