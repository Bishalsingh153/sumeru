package web_test

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"sumeru/core/server/web"
)

func TestFlashFromQueryMessage_truncatesLongErrorBody(t *testing.T) {
	long := strings.Repeat("x", 600)
	flash, ok := web.FlashFromQueryMessage("error:" + long)
	if !ok {
		t.Fatal("expected flash")
	}
	if utf8.RuneCountInString(flash.Body) > 500 {
		t.Fatalf("body runes=%d want <=500", utf8.RuneCountInString(flash.Body))
	}
}

func TestSignedUIDCookie_rejectsExpired(t *testing.T) {
	val := web.MintSignedUIDCookieForTest(7, -time.Second)
	if uid, ok := web.ParseSignedUIDCookieForTest(val); ok || uid != 0 {
		t.Fatalf("expected expired cookie rejected, uid=%d ok=%v", uid, ok)
	}
}
