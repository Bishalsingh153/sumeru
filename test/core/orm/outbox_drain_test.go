package orm_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"sumeru/core/orm"
)

func registerOutboxStub(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	orm.SetDBForTest(orm.NewDBWrapper(db))
	t.Cleanup(func() {
		_ = db.Close()
		orm.ResetDBForTest()
	})
	orm.RegisterStubModelForTest(t, "sys.outbox.event", []orm.FieldDefinition{
		{Name: "name", Type: orm.Char},
	})
	return mock
}

func TestDrainOutboxOnce_idleQueueSkipsDrain(t *testing.T) {
	mock := registerOutboxStub(t)
	mock.ExpectQuery(`SELECT id, name, COALESCE\(payload_json`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "payload_json", "actor"}))

	if n := orm.DrainOutboxOnce(context.Background()); n != 0 {
		t.Fatalf("published = %d, want 0", n)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
