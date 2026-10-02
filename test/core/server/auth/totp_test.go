package auth_test

import (
	"testing"

	"sumeru/core/server/auth"
)

func TestValidateTOTP_rejectsEmptySecret(t *testing.T) {
	if auth.ValidateTOTP("", "123456") {
		t.Fatal("empty secret should fail")
	}
}

func TestValidateTOTP_window(t *testing.T) {
	secret, err := auth.GenerateTOTPSecret()
	if err != nil {
		t.Fatal(err)
	}
	if auth.ValidateTOTP(secret, "abc") {
		t.Fatal("non-numeric code should fail")
	}
	if auth.ValidateTOTP(secret, "000000") && auth.ValidateTOTP(secret, "111111") {
		t.Fatal("random secret should not match two arbitrary codes")
	}
}
