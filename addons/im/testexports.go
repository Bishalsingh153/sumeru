package im

import "context"

// UserDisplayNameForTest exposes peer display name resolution for tests.
func UserDisplayNameForTest(ctx context.Context, userID int) string {
	return UserDisplayName(ctx, userID)
}
