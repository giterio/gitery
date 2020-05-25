package jwt

import (
	"testing"
	"time"
)

type JwtPayload struct {
	Iss string `json:"iss,omitempty"` // issuer
	Exp int64  `json:"exp,omitempty"` // expiration time
	Sub string `json:"sub,omitempty"` // subject
	Aud string `json:"aud,omitempty"` // audience
	Nbf int64  `json:"nbf,omitempty"` // Not Before
	Iat int64  `json:"iat"`           // Issued At
	Jti int64  `json:"jti,omitempty"` // JWT ID
	Pub struct {
		ID    *int   `json:"userID"`
		Email string `json:"email"`
	} `json:"pub"`
}

func TestEncode(t *testing.T) {
	secret := "this is screct"
	id := 1
	payload := JwtPayload{
		Sub: "123",
		Exp: time.Now().Unix() + 100000,
		Pub: struct {
			ID    *int   `json:"userID"`
			Email string `json:"email"`
		}{
			ID:    &id,
			Email: "Murphy@jwt.com",
		},
	}
	token, err := Encode(payload, secret)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	payloadDecoded := JwtPayload{}
	err = Decode(token, secret, payloadDecoded)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if payload.Exp != payloadDecoded.Exp {
		t.Errorf("Error: payload data not matched")
	}
	if payloadDecoded.Pub.Email != payload.Pub.Email || *payloadDecoded.Pub.ID != *payload.Pub.ID {
		t.Errorf("Error: payload pub not matched")
	}
}
