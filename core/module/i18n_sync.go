package module

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sumeru/core/orm"
)

// SyncAddonTranslationPOFiles imports gettext PO files from addonPath/i18n/*.po.
func SyncAddonTranslationPOFiles(ctx context.Context, addonPath, moduleName string) error {
	if orm.DB == nil {
		return nil
	}
	dir := filepath.Join(addonPath, "i18n")
	matches, err := filepath.Glob(filepath.Join(dir, "*.po"))
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return nil
	}
	var firstErr error
	importedTotal := 0
	for _, poPath := range matches {
		f, err := os.Open(poPath)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		langHint, entries, err := orm.ParsePO(f)
		_ = f.Close()
		if err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("%s: %w", filepath.Base(poPath), err)
			}
			continue
		}
		lang := strings.TrimSpace(langHint)
		if lang == "" {
			lang = orm.LangCodeFromPOFilename(poPath)
		}
		n, err := orm.ImportTranslationsPO(ctx, orm.DB, moduleName, lang, entries)
		if err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("%s: %w", filepath.Base(poPath), err)
			}
			continue
		}
		importedTotal += n
	}
	if firstErr != nil {
		return firstErr
	}
	_ = importedTotal
	return nil
}
