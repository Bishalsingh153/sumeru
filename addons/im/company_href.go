package im

import (
	"context"
	"fmt"

	"sumeru/core/orm"
)

// CompanySettingsFormHref returns a workspace URL to edit company IM settings when the user may read core.company.
func CompanySettingsFormHref(ctx context.Context, companyID int) string {
	if companyID <= 0 {
		return ""
	}
	uid := orm.SecurityUID(ctx)
	if uid <= 0 {
		return ""
	}
	if err := orm.CheckModelAccess(ctx, uid, "core.company", "read"); err != nil {
		return ""
	}
	return fmt.Sprintf("/web?model=core.company&view_type=form&id=%d", companyID)
}
