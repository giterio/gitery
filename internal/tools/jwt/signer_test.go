package jwt

import "testing"

func TestSign(t *testing.T) {
	secret := "this is screct"
	msg := "this is message"
	sig := sign(msg, secret)
	if ok := verify(msg, sig, secret); !ok {
		t.Errorf("can not verify the sig")
	}
}
