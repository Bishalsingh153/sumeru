package web_test

import (
	"testing"
	"time"

	"sumeru/core/server/web"
)

func TestSignedUIDCookie_roundTrip(t *testing.T) {
	val := web.MintSignedUIDCookieForTest(42, time.Minute)
	uid, ok := web.ParseSignedUIDCookieForTest(val)
	if !ok || uid != 42 {
		t.Fatalf("parse signed uid: ok=%v uid=%d", ok, uid)
	}
}

func TestSignedUIDCookie_rejectsTamper(t *testing.T) {
	val := web.MintSignedUIDCookieForTest(42, time.Minute)
	tampered := val[:len(val)-2] + "xx"
	if uid, ok := web.ParseSignedUIDCookieForTest(tampered); ok || uid != 0 {
		t.Fatalf("expected tampered cookie rejected, got uid=%d ok=%v", uid, ok)
	}
}
