package orm

import (
	"context"
	"fmt"
)

// ensureExtraIndexes creates composite indexes not expressible via single-field Index flags.
func ensureExtraIndexes(ctx context.Context) error {
	if DB == nil {
		return nil
	}
	return ensureSysTranslationUniqueIndex(ctx)
}

func ensureSysTranslationUniqueIndex(ctx context.Context) error {
	tablePhysical := MustModelToTableName("sys.translation")
	if tablePhysical == "" {
		return nil
	}
	ok, err := tableExists(ctx, tablePhysical)
	if err != nil || !ok {
		return err
	}
	tableQuoted := MustQuotedTableName("sys.translation")
	langCol, err := QuotedColumnForModel("sys.translation", "lang")
	if err != nil {
		return err
	}
	srcCol, err := QuotedColumnForModel("sys.translation", "src")
	if err != nil {
		return err
	}
	moduleCol, err := QuotedColumnForModel("sys.translation", "module")
	if err != nil {
		return err
	}
	idxName := "sys_translation_lang_src_module_uidx"
	q := fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS %s ON %s (%s, %s, %s)",
		quoteIdent(idxName), tableQuoted, langCol, srcCol, moduleCol)
	_, err = DB.ExecContext(ctx, q)
	return err
}
