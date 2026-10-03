package orm_test

import (
	"strings"
	"testing"

	"sumeru/core/orm"
)

func TestParsePO_simple(t *testing.T) {
	raw := `# comment
msgid "Hello"
msgstr "Bonjour"

msgctxt "base"
msgid "Settings"
msgstr "Paramètres"
`
	lang, entries, err := orm.ParsePO(strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	if lang != "" {
		t.Fatalf("lang hint = %q", lang)
	}
	if len(entries) != 2 {
		t.Fatalf("entries = %d", len(entries))
	}
	if entries[0].MsgID != "Hello" || entries[0].MsgStr != "Bonjour" {
		t.Fatalf("first entry: %+v", entries[0])
	}
	if entries[1].Context != "base" || entries[1].MsgID != "Settings" {
		t.Fatalf("second entry: %+v", entries[1])
	}
}

func TestParsePO_headerLanguage(t *testing.T) {
	raw := `msgid ""
msgstr ""
"Language: fr_FR\n"

msgid "Hi"
msgstr "Salut"
`
	lang, entries, err := orm.ParsePO(strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	if lang != "fr_FR" {
		t.Fatalf("lang = %q", lang)
	}
	if len(entries) != 1 || entries[0].MsgStr != "Salut" {
		t.Fatalf("entries: %+v", entries)
	}
}

func TestNormalizeLangCode(t *testing.T) {
	if orm.NormalizeLangCode("fr") != "fr_FR" {
		t.Fatal("fr")
	}
	if orm.NormalizeLangCode("en_US") != "en_US" {
		t.Fatal("en_US")
	}
}

func TestLangCodeFromPOFilename(t *testing.T) {
	if got := orm.LangCodeFromPOFilename("/tmp/es.po"); got != "es_ES" {
		t.Fatalf("got %q", got)
	}
}
