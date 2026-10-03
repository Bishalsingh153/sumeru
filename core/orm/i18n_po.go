package orm

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// POEntry is one translated string from a gettext PO catalog.
type POEntry struct {
	Context string
	MsgID   string
	MsgStr  string
}

var poLanguageHeader = regexp.MustCompile(`(?i)Language:\s*([^\s\\]+)`)

// ParsePO reads a gettext PO file into entries. langHint is from the header Language: line when present.
func ParsePO(r io.Reader) (langHint string, entries []POEntry, err error) {
	sc := bufio.NewScanner(r)
	var cur POEntry
	var collecting string
	var target *string

	flush := func() {
		if cur.MsgID == "" && cur.MsgStr == "" {
			return
		}
		if cur.MsgID == "" {
			if m := poLanguageHeader.FindStringSubmatch(cur.MsgStr); len(m) == 2 {
				langHint = strings.TrimSpace(m[1])
			}
			cur = POEntry{}
			return
		}
		entries = append(entries, cur)
		cur = POEntry{}
	}

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			if collecting != "" {
				flush()
				collecting = ""
				target = nil
			}
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "msgctxt ") {
			cur.Context = unquotePO(line[len("msgctxt "):])
			continue
		}
		if strings.HasPrefix(line, "msgid ") {
			if collecting != "" && cur.MsgID != "" {
				flush()
			}
			collecting = "msgid"
			cur.MsgID = unquotePO(line[len("msgid "):])
			target = &cur.MsgID
			continue
		}
		if strings.HasPrefix(line, "msgstr ") {
			collecting = "msgstr"
			cur.MsgStr = unquotePO(line[len("msgstr "):])
			target = &cur.MsgStr
			continue
		}
		if strings.HasPrefix(line, `"`) && target != nil {
			*target += unquotePO(line)
		}
	}
	if err := sc.Err(); err != nil {
		return "", nil, err
	}
	if collecting != "" {
		flush()
	}
	return langHint, entries, nil
}

func unquotePO(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		body := s[1 : len(s)-1]
		return strings.ReplaceAll(body, `\"`, `"`)
	}
	return s
}

// NormalizeLangCode maps PO language tags to core.lang-style codes (e.g. fr → fr_FR).
func NormalizeLangCode(code string) string {
	code = strings.TrimSpace(code)
	code = strings.ReplaceAll(code, "-", "_")
	if code == "" {
		return ""
	}
	if strings.Contains(code, "_") {
		return code
	}
	switch code {
	case "en":
		return "en_US"
	default:
		return code + "_" + strings.ToUpper(code)
	}
}

// LangCodeFromPOFilename extracts a language code from names like fr_FR.po or es.po.
func LangCodeFromPOFilename(name string) string {
	base := name
	if i := strings.LastIndex(base, "/"); i >= 0 {
		base = base[i+1:]
	}
	base = strings.TrimSuffix(base, ".po")
	return NormalizeLangCode(base)
}

// ImportTranslationsPO upserts PO entries for moduleName and lang.
func ImportTranslationsPO(ctx context.Context, db DBWrapper, moduleName, lang string, entries []POEntry) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("db required")
	}
	lang = NormalizeLangCode(lang)
	if lang == "" {
		return 0, fmt.Errorf("language code required")
	}
	mod := strings.TrimSpace(moduleName)
	imported := 0
	for _, e := range entries {
		src := strings.TrimSpace(e.MsgID)
		if src == "" {
			continue
		}
		if err := upsertTranslation(ctx, db, lang, src, e.MsgStr, mod); err != nil {
			return imported, err
		}
		imported++
	}
	return imported, nil
}

// ExportTranslationsPO writes PO catalog for moduleName and lang.
func ExportTranslationsPO(ctx context.Context, db DBWrapper, moduleName, lang string, w io.Writer) (int, error) {
	lang = NormalizeLangCode(lang)
	mod := strings.TrimSpace(moduleName)
	tbl := MustQuotedTableName("sys.translation")
	rows, err := db.QueryContext(ctx,
		`SELECT src, value FROM `+tbl+` WHERE lang = $1 AND module = $2 ORDER BY src`,
		lang, mod)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if _, err := fmt.Fprintf(w, "# Sumeru translation catalog\nmsgid \"\"\nmsgstr \"\"\n\"Language: %s\\n\"\n\n", lang); err != nil {
		return 0, err
	}

	count := 0
	for rows.Next() {
		var src, value string
		if err := rows.Scan(&src, &value); err != nil {
			return count, err
		}
		if err := writePOEntry(w, "msgid", src); err != nil {
			return count, err
		}
		if err := writePOEntry(w, "msgstr", value); err != nil {
			return count, err
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return count, err
		}
		count++
	}
	return count, rows.Err()
}

func writePOEntry(w io.Writer, keyword, s string) error {
	esc := func(x string) string {
		x = strings.ReplaceAll(x, `\`, `\\`)
		return strings.ReplaceAll(x, `"`, `\"`)
	}
	if !strings.Contains(s, "\n") {
		_, err := fmt.Fprintf(w, `%s "%s"`+"\n", keyword, esc(s))
		return err
	}
	parts := strings.Split(s, "\n")
	if _, err := fmt.Fprintf(w, `%s ""`+"\n", keyword); err != nil {
		return err
	}
	for i, p := range parts {
		line := `"` + esc(p) + `"`
		if i < len(parts)-1 {
			line += `\n`
		}
		if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
			return err
		}
	}
	return nil
}
