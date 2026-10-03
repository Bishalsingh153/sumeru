# Internal chat (`im`)

P2P messaging between internal users via the shell activity dock **Messages** tab.

- Model: `im.message`
- No record chatter, email, activities, followers, or templates
- Company flag: `im_enabled` on `core.company` (falls back to legacy `mail_chatter_enabled` until migrated)
- **Log** tab in the activity dock remains `sys.audit` (from the **audit** addon)

Upgrade from the old Discuss addon: install **im**, uninstall **mail**, and migrate direct message rows if needed.
