package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

func base64URLMarshal(raw interface{}) (base64Str string, err error) {
	bytes, err := json.Marshal(raw)
	if err != nil {
		return
	}
	base64Str = base64.URLEncoding.EncodeToString(bytes)
	return
}

func base64URLUnmarshal(base64Str string, v interface{}) (err error) {
	bytes, err := base64.URLEncoding.DecodeString(base64Str)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, v)
	return
}

func sign(src string, secret string) (sig string) {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(src))
	sig = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}

func verify(msg string, sig string, secret string) bool {
	return sig == sign(msg, secret)
}
