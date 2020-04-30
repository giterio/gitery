package jwt

import (
	"errors"
	"strings"
	"time"
)

// Payload ...
type Payload struct {
	Iss string      `json:"iss,omitempty"` // issuer
	Exp int64       `json:"exp,omitempty"` // expiration time
	Sub string      `json:"sub,omitempty"` // subject
	Aud string      `json:"aud,omitempty"` // audience
	Nbf int64       `json:"nbf,omitempty"` // Not Before
	Iat int64       `json:"iat,omitempty"` // Issued At
	Jti int64       `json:"jti,omitempty"` // JWT ID
	Pub interface{} `json:"pub,omitempty"` // Public message
}

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// Encode ...
func Encode(payload Payload, secret string) (token string, err error) {
	header := header{
		Alg: "HS256",
		Typ: "JWT",
	}
	headerBase64URL, err := base64URLMarshal(header)
	if err != nil {
		return
	}

	payloadBase64URL, err := base64URLMarshal(payload)
	if err != nil {
		return
	}

	msg := headerBase64URL + "." + payloadBase64URL
	sig := sign(msg, secret)
	token = msg + "." + sig
	return
}

// Decode ...
func Decode(token string, secret string) (payload Payload, err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		err = errors.New("Invalid token: wrong format")
		return
	}
	payload = Payload{}
	err = base64URLUnmarshal(parts[1], &payload)
	if err != nil {
		err = errors.New("Invalid token: payload not parsable")
		return
	}
	if payload.Exp != 0 && time.Now().Unix() > payload.Exp {
		err = errors.New("Invalid token: token expired")
		return
	}
	msg := parts[0] + "." + parts[1]
	if ok := verify(msg, parts[2], secret); !ok {
		err = errors.New("Invalid token: verification failed")
		return
	}
	return
}
