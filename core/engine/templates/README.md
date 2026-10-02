# Engine HTML templates

Root path: `core/engine/templates` (`templates_path` in config). Relative paths are defined in [`render/templates_paths.go`](../render/templates_paths.go).

| Folder | Role | Loaded by |
|--------|------|-----------|
| `shell/` | Authenticated app chrome (`base.html`) | `render.RenderPage` |
| `partials/` | Shared `{{define}}` fragments (menu icons, home hub, apps forms, login brand, auth assets) | Parsed with shell layout, shell inner pages, or auth pages |
| `pages/` | Main column HTML injected into `base.html` | `web.renderShellPage` / `executeInnerTemplate` |
| `auth/` | Standalone login, TOTP, setup wizards | `web.getAuthTemplate`, setup handler |
| `portal/` | Portal home and share pages | `web.loadPortalTemplates` |

## Parse sets

- **Shell layout:** `shell/base.html` + `partials/menu_icons.html`
- **Shell inner page:** `pages/<page>.html` + menu, home hub, and apps partials
- **Auth page:** `auth/<page>.html` + `partials/login_brand_aside.html` + `partials/auth_head.html`
