package web_test

import (
	"testing"

	"sumeru/core/server/web"
)

func TestPartitionListSections(t *testing.T) {
	rows := []map[string]interface{}{
		{"id": 1, "state": "draft", "name": "A"},
		{"id": 2, "state": "done", "name": "B"},
		{"id": 3, "state": "draft", "name": "C"},
	}
	sections := web.PartitionListSectionsForTest(rows, "state")
	if len(sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(sections))
	}
	if sections[0].Count+sections[1].Count != 3 {
		t.Fatalf("row count mismatch")
	}
}

func TestPartitionListSections_emptyGroupLabel(t *testing.T) {
	rows := []map[string]interface{}{
		{"id": 1, "mobile": nil, "name": "A"},
		{"id": 2, "mobile": "", "name": "B"},
	}
	sections := web.PartitionListSectionsForTest(rows, "mobile")
	if len(sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(sections))
	}
	if sections[0].Label != "(Empty)" {
		t.Fatalf("label = %q, want (Empty)", sections[0].Label)
	}
	if sections[0].Count != 2 {
		t.Fatalf("count = %d", sections[0].Count)
	}
}
