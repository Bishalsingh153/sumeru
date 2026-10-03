package mail_test

import (
	"context"
	"testing"

	"sumeru/addons/mail"
)

func TestListAttachmentsForRecordInvalid(t *testing.T) {
	_, err := mail.ListAttachmentsForRecord(context.Background(), "", 0, 10)
	if err == nil {
		t.Fatal("expected error")
	}
	_, err = mail.ListAttachmentsForRecord(context.Background(), "not.a.model", 1, 10)
	if err == nil {
		t.Fatal("expected unknown model error")
	}
}
