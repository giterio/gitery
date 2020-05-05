package jwt

import (
	"errors"
	"strings"
)

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// Encode ...
func Encode(payload interface{}, secret string) (token string, err error) {
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
	// encrypt URLBase64(header).URLBase64(payload) to generate signature
	sig := sign(msg, secret)
	token = msg + "." + sig
	return
}

// Decode ...
func Decode(token string, secret string, payload interface{}) (err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		err = errors.New("Invalid token: wrong format")
		return
	}
	err = base64URLUnmarshal(parts[1], payload)
	if err != nil {
		err = errors.New("Invalid token: payload not parsable")
		return
	}
	msg := parts[0] + "." + parts[1]
	if ok := verify(msg, parts[2], secret); !ok {
		err = errors.New("Invalid token: verification failed")
		return
	}
	return
}
