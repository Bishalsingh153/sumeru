package swcmeta_test

import (
	"context"
	"testing"

	"sumeru/core/engine/parser"
	"sumeru/core/engine/swcmeta"
	"sumeru/core/report"
)

func TestApplyReportMasterDataDefault(t *testing.T) {
	arch := swcmeta.ViewArch{Type: "list", Model: "core.company"}
	swcmeta.ApplyReportToArchForTest(context.Background(), &arch, "core.company", nil, nil)
	if arch.Report == nil || !arch.Report.Download || !arch.Report.Upload {
		t.Fatalf("expected master-data default report, got %+v", arch.Report)
	}
}

func TestApplyReportSkipsNonMasterModel(t *testing.T) {
	arch := swcmeta.ViewArch{Type: "list", Model: "sys.audit"}
	swcmeta.ApplyReportToArchForTest(context.Background(), &arch, "sys.audit", nil, nil)
	if arch.Report != nil {
		t.Fatalf("expected no report for sys.audit, got %+v", arch.Report)
	}
}

func TestApplyReportInvisibleOnActiveView(t *testing.T) {
	active, err := parser.ParseViewFromArch(`<list><report invisible="1" download="csv"/></list>`)
	if err != nil {
		t.Fatal(err)
	}
	arch := swcmeta.ViewArch{Type: "list", Model: "core.user"}
	swcmeta.ApplyReportToArchForTest(context.Background(), &arch, "core.user", active, []string{"list", "form"})
	if arch.Report != nil {
		t.Fatalf("expected report hidden on invisible list view, got %+v", arch.Report)
	}
}

func TestUnionReportFromFormCapabilities(t *testing.T) {
	form, err := parser.ParseViewFromArch(`<form><report download="csv,pdf" upload="bulk"/><field name="name"/></form>`)
	if err != nil {
		t.Fatal(err)
	}
	formCaps := report.CapabilitiesFromView(form)
	if !formCaps.HasDownload() || !formCaps.BulkUpload {
		t.Fatalf("unexpected form caps: %+v", formCaps)
	}
	list, err := parser.ParseViewFromArch(`<list><field name="name"/></list>`)
	if err != nil {
		t.Fatal(err)
	}
	listCaps := report.CapabilitiesFromView(list)
	merged := swcmeta.UnionReportCapabilitiesForTest(listCaps, formCaps)
	if !merged.HasDownload() || !merged.BulkUpload {
		t.Fatalf("expected merged caps from form, got %+v", merged)
	}
}

func TestCapabilitiesFromViewRespectsReportInvisible(t *testing.T) {
	v, err := parser.ParseViewFromArch(`<list><report invisible="1" download="csv"/></list>`)
	if err != nil {
		t.Fatal(err)
	}
	caps := report.CapabilitiesFromView(v)
	if caps.HasDownload() || caps.BulkUpload {
		t.Fatalf("invisible report should yield empty caps, got %+v", caps)
	}
}
