package im_test

import (
	"context"
	"testing"

	"sumeru/addons/im"
	"sumeru/core/orm"
)

func TestSortDirectMessages(t *testing.T) {
	rows := []map[string]interface{}{
		{"create_date": "2024-01-02T00:00:00Z"},
		{"create_date": "2024-01-01T00:00:00Z"},
	}
	im.SortDirectMessagesForTest(rows)
	if orm.AsString(rows[0]["create_date"]) != "2024-01-01T00:00:00Z" {
		t.Fatal("expected ascending sort")
	}
}

func TestUserDisplayName_noDB(t *testing.T) {
	name := im.UserDisplayNameForTest(context.Background(), 99)
	if name == "" {
		t.Fatal("expected fallback name")
	}
}

func TestIMEnabled_defaultTrueWithoutDB(t *testing.T) {
	if !im.IMEnabled(context.Background()) {
		t.Fatal("default IM enabled")
	}
}

func TestPostDirectMessage_requiresBody(t *testing.T) {
	_, err := im.PostDirectMessage(context.Background(), 1, 2, "   ")
	if err == nil {
		t.Fatal("expected error for empty body")
	}
	_, err = im.PostDirectMessage(context.Background(), 1, 1, "hi")
	if err == nil {
		t.Fatal("expected error for self peer")
	}
}

func TestSearchInternalUsers_nilDB(t *testing.T) {
	users, err := im.SearchInternalUsers(context.Background(), 1, "john doe", 0)
	if err != nil || len(users) != 0 {
		t.Fatalf("nil db: %v %v", users, err)
	}
}

func TestDirectHelpers_nilDB(t *testing.T) {
	ctx := context.Background()
	users, err := im.SearchInternalUsers(ctx, 1, "a", 5)
	if err != nil || len(users) != 0 {
		t.Fatalf("SearchInternalUsers nil db: %v %v", users, err)
	}
	convs, err := im.ListDirectConversations(ctx, 1, 5)
	if err != nil || len(convs) != 0 {
		t.Fatalf("ListDirectConversations nil db: %v %v", convs, err)
	}
	thread, err := im.ListDirectThread(ctx, 1, 2)
	if err != nil || len(thread) != 0 {
		t.Fatalf("ListDirectThread nil db: %v %v", thread, err)
	}
	_, err = im.ListDirectThread(ctx, 1, 1)
	if err == nil {
		t.Fatal("expected invalid peer")
	}
	_, err = im.PostDirectMessage(ctx, 1, 2, "hello")
	if err == nil {
		t.Fatal("expected error without message model registered")
	}
}
