package im

import (
	"context"
	"database/sql"

	"sumeru/core/orm"
)

// IMEnabled reports whether internal chat is enabled for the active company (default true).
func IMEnabled(ctx context.Context) bool {
	if orm.DB == nil {
		return true
	}
	companyID := orm.CompanyIDFromContext(ctx)
	if companyID <= 0 {
		companyID = orm.ActiveCompanyIDForUser(ctx, orm.SecurityUID(ctx))
	}
	tn := orm.MustQuotedTableName("core.company")
	var legacyChatter sql.NullBool
	q := `SELECT mail_chatter_enabled FROM ` + tn
	if companyID > 0 {
		q += ` WHERE id = $1`
		err := orm.DB.QueryRowContext(ctx, q, companyID).Scan(&legacyChatter)
		if err != nil {
			return true
		}
	} else {
		q += ` ORDER BY id ASC LIMIT 1`
		if err := orm.DB.QueryRowContext(ctx, q).Scan(&legacyChatter); err != nil {
			return true
		}
	}
	if legacyChatter.Valid {
		return legacyChatter.Bool
	}
	return true
}
