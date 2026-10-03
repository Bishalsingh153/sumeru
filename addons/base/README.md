# Base

Sumeru kernel module: users, companies, partners, security administration, and settings hub.

Related auto-installed modules:

| Module | Role |
|--------|------|
| **audit** | Audit trail, retention, exports |
| **platform** | Attachments, sequences, parameters, reports, import templates |
| **geo** | Regional data (countries, states, cities, currency) |
| **i18n** | Internationalization (languages, translations, PO catalogs) |
| **im** | Internal user-to-user chat (activity dock Messages tab) |

See [addons/README.md](../README.md) for the standard addon layout.

## Install

Base is installed automatically with the server. To reload XML after changes:

```bash
go run ./cmd/sumeru -- -c sumeru.conf -u base -u audit -u platform -u geo -u i18n -u im
```

PO catalogs under `addons/base/i18n/*.po` load automatically on `-u base`.
