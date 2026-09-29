package swcmeta

import (
	"context"
	"reflect"

	"sumeru/core/engine/parser"
	"sumeru/core/report"
)

func SerializeSheetForTest(ctx context.Context, model string, s *parser.Sheet) *ArchSheet {
	return serializeSheet(ctx, model, s)
}

func SerializeGroupForTest(ctx context.Context, model string, g parser.Group) ArchGroup {
	return serializeGroup(ctx, model, g)
}

func SerializeDivForTest(ctx context.Context, model string, d parser.Div) ArchDiv {
	return serializeDiv(ctx, model, d)
}

func FormMetaForModelForTest(model string) *FormMeta { return formMetaForModel(model) }

func SerializeFieldsForTest(ctx context.Context, fields []parser.Field) []ArchField {
	return serializeFields(ctx, fields)
}

func EnrichFieldForTest(model string, f ArchField) ArchField { return enrichField(model, f) }

func WorkspacePayloadTypeForTest() reflect.Type { return reflect.TypeOf(WorkspacePayload{}) }

func ApplyReportToArchForTest(
	ctx context.Context,
	arch *ViewArch,
	model string,
	activeView *parser.View,
	actionViewModes []string,
) {
	ApplyReportToArch(ctx, arch, model, activeView, actionViewModes)
}

func UnionReportCapabilitiesForTest(a, b report.Capabilities) report.Capabilities {
	return unionReportCapabilities(a, b)
}
