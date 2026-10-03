# Internal chat (`im`)

P2P messaging between internal users via the shell activity dock **Messages** tab.

- Model: `im.message`
- **Messages** = direct chat between internal users (not tied to the record you have open).
- **Log** = change history for the open record (`sys.audit` via the **audit** addon), visible when you can read that record.
- No record chatter, email, activities, followers, or templates
- Company flag: `im_enabled` on `core.company` (falls back to legacy `mail_chatter_enabled` when present)

- **Unread** count on the Messages tab; messages marked read when you open a thread (`is_read`).
- Optional **record link** on each message (`res_model` / `res_id`) when sent while a form record is open.
- **`GET /web/swc/direct/status`**: `enabled`, `unread`, `companyFormHref` for the activity UI.

Message attachments are `sys.attachment` rows served at `/web/content/<id>` with ACL checks. Optional virus scanning: `orm.SetAttachmentScanner` (see platform README / `core/orm/attachment_binary.go`).

Upgrade from the old Discuss addon: install **im**, uninstall **mail**, and migrate direct message rows if needed. After model changes, run **`-u im`** so `is_read`, `res_model`, and `res_id` columns exist.
