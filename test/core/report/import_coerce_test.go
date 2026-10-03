package report_test

import (
	"context"
	"testing"

	"sumeru/core/orm"
	"sumeru/core/report"
)

type stubImportModel struct{}

func (stubImportModel) ModelName() string { return "stub.import" }
func (stubImportModel) Fields() []orm.FieldDefinition {
	return []orm.FieldDefinition{
		{Name: "active", Type: orm.Boolean},
		{Name: "amount", Type: orm.Float},
		{Name: "qty", Type: orm.Integer},
	}
}

func TestCoerceImportBooleanAndFloat(t *testing.T) {
	out, errs := report.CoerceImportValuesForTest(context.Background(), stubImportModel{}, map[string]interface{}{
		"active": "true",
		"amount": "3.5",
	})
	if len(errs) > 0 {
		t.Fatalf("errs=%v", errs)
	}
	if out["active"] != true {
		t.Fatalf("active=%v", out["active"])
	}
	if out["amount"] != 3.5 {
		t.Fatalf("amount=%v", out["amount"])
	}
}

func TestParseUploadContentInvalidXLSX(t *testing.T) {
	_, err := report.ParseUploadContentForTest("bad.xlsx", []byte("not-a-zip"))
	if err == nil {
		t.Fatal("expected xlsx error")
	}
}

func TestParseUploadContentCSVPassThrough(t *testing.T) {
	in := []byte("a,b\n1,2\n")
	out, err := report.ParseUploadContentForTest("data.csv", in)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != string(in) {
		t.Fatalf("got %q", out)
	}
}

type selImportModel struct{}

func (selImportModel) ModelName() string { return "sel" }
func (selImportModel) Fields() []orm.FieldDefinition {
	return []orm.FieldDefinition{{Name: "state", Type: orm.Selection, Selection: [][]string{{"a", "A"}}}}
}

type dateImportModel struct{}

func (dateImportModel) ModelName() string { return "d" }
func (dateImportModel) Fields() []orm.FieldDefinition {
	return []orm.FieldDefinition{{Name: "day", Type: orm.Date}}
}

func TestCoerceImportDateParsing(t *testing.T) {
	out, errs := report.CoerceImportValuesForTest(context.Background(), dateImportModel{}, map[string]interface{}{"day": "2024-01-02"})
	if len(errs) > 0 {
		t.Fatalf("errs=%v", errs)
	}
	if out["day"] == nil {
		t.Fatal("expected date")
	}
}

func TestCoerceImportSelectionInvalid(t *testing.T) {
	_, errs := report.CoerceImportValuesForTest(context.Background(), selImportModel{}, map[string]interface{}{"state": "bad"})
	if len(errs) == 0 {
		t.Fatal("expected selection error")
	}
}

func TestDryRunBulkImportMissingBatch(t *testing.T) {
	_, _, err := report.DryRunBulkImport(context.Background(), report.PreviewBulkImportInput{BatchID: 1})
	if err == nil {
		t.Fatal("expected error without database batch")
	}
}

func TestRunBulkJobCronNoDB(t *testing.T) {
	if err := report.RunBulkJobCronForTest(context.Background()); err != nil {
		t.Fatalf("cron without db: %v", err)
	}
}

func TestQueueExportJobMissingModel(t *testing.T) {
	_, err := report.QueueExportJob(context.Background(), "no.such.model", []string{"name"}, nil, "csv", 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestApplyImportTemplateMappingNoDB(t *testing.T) {
	if m := report.ApplyImportTemplateMapping(context.Background(), "core.partner", []string{"Name"}, 1); m != nil {
		t.Fatalf("expected nil without db, got %v", m)
	}
}

func TestSaveImportTemplateNoRegistry(t *testing.T) {
	_, err := report.SaveImportTemplate(context.Background(), "T", "core.partner", "create", "[]", "{}", 1, false)
	if err == nil {
		t.Fatal("expected error without template model")
	}
}

func TestCountBatchRowsMissing(t *testing.T) {
	_, err := report.CountBatchRows(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFetchRowsUnlimitedMissingModel(t *testing.T) {
	_, err := report.FetchRowsUnlimitedForTest(context.Background(), "missing.model", nil, 0)
	if err == nil {
		t.Fatal("expected error")
	}
}

type dtImportModel struct{}

func (dtImportModel) ModelName() string { return "dt" }
func (dtImportModel) Fields() []orm.FieldDefinition {
	return []orm.FieldDefinition{{Name: "at", Type: orm.DateTime}}
}

func TestCoerceImportDateTimeRFC3339(t *testing.T) {
	out, errs := report.CoerceImportValuesForTest(context.Background(), dtImportModel{}, map[string]interface{}{"at": "2024-01-02T15:04:05Z"})
	if len(errs) > 0 {
		t.Fatalf("errs=%v", errs)
	}
	if out["at"] == nil {
		t.Fatal("expected time")
	}
}

func TestCoerceImportIntegerInvalid(t *testing.T) {
	_, errs := report.CoerceImportValuesForTest(context.Background(), stubImportModel{}, map[string]interface{}{"qty": "not-int"})
	if len(errs) == 0 {
		t.Fatal("expected int error")
	}
}

func TestExportCSVUnlimitedMissingModel(t *testing.T) {
	_, err := report.ExportCSVUnlimitedForTest(context.Background(), report.ExportCSVInput{Model: "nope", Fields: []string{"x"}})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRunExportJobMissingModel(t *testing.T) {
	err := report.RunExportJobForTest(context.Background(), map[string]interface{}{"target_model": ""})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPreviewBulkImportZeroBatch(t *testing.T) {
	_, err := report.PreviewBulkImport(context.Background(), report.PreviewBulkImportInput{})
	if err == nil {
		t.Fatal("expected error")
	}
}
