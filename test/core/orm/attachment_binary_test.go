package orm_test

import (
	"context"
	"path/filepath"
	"testing"

	"sumeru/core/orm"
)

func TestCreateBinaryAttachmentRequiresModel(t *testing.T) {
	if err := orm.InitFilestore(filepath.Join(t.TempDir(), "fs")); err != nil {
		t.Fatal(err)
	}
	_, err := orm.CreateBinaryAttachment(context.Background(), orm.CreateBinaryAttachmentInput{
		Name: "x.txt",
		Data: []byte("data"),
	})
	if err == nil {
		t.Fatal("expected error without sys.attachment registry")
	}
}

func TestSetAttachmentScannerNil(t *testing.T) {
	orm.SetAttachmentScanner(nil)
}
