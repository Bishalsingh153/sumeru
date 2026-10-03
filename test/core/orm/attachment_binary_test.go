package orm_test

import (
	"context"
	"io"
	"path/filepath"
	"strings"
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

func TestAttachmentScannerRunsBeforeCreate(t *testing.T) {
	if err := orm.InitFilestore(filepath.Join(t.TempDir(), "fs")); err != nil {
		t.Fatal(err)
	}
	called := false
	orm.SetAttachmentScanner(func(_ context.Context, _, _ string, _ io.Reader) error {
		called = true
		return nil
	})
	defer orm.SetAttachmentScanner(nil)
	_, _ = orm.CreateBinaryAttachment(context.Background(), orm.CreateBinaryAttachmentInput{
		Name: "note.txt",
		Data: []byte("hello"),
	})
	if !called {
		t.Fatal("expected scanner hook to run")
	}
}

func TestAttachmentScannerBlocksUpload(t *testing.T) {
	if err := orm.InitFilestore(filepath.Join(t.TempDir(), "fs")); err != nil {
		t.Fatal(err)
	}
	orm.SetAttachmentScanner(func(_ context.Context, _, _ string, _ io.Reader) error {
		return context.Canceled
	})
	defer orm.SetAttachmentScanner(nil)
	_, err := orm.CreateBinaryAttachment(context.Background(), orm.CreateBinaryAttachmentInput{
		Name: "bad.bin",
		Data: []byte("x"),
	})
	if err == nil || !strings.Contains(err.Error(), "scan") {
		t.Fatalf("expected scan error, got %v", err)
	}
}
