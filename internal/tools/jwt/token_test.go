package jwt

import (
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	secret := "this is screct"
	payload := Payload{
		Sub: "123",
		Exp: time.Now().Unix() + 100000,
		Pub: map[string]interface{}{
			"Name":  "Murphy",
			"Email": "Murphy@jwt.com",
		},
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
	pub1, ok1 := payload.Pub.(map[string]interface{})
	pub2, ok2 := payloadDecoded.Pub.(map[string]interface{})
	if !ok1 || !ok2 || pub1["Name"] != pub2["Name"] {
		t.Errorf("Error: payload pub not matched")
	}
}
