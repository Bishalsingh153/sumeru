# Platform addon

Bulk import/export, import templates, sequences, and shared attachment storage patterns.

## Bulk jobs

- Model: `sys.bulk.import` (import and export batches).
- Cron: **`platform.bulk_job`** must be active for queued imports (>500 rows) and background exports.
- List toolbar: when a list has more than 500 rows, use **Export all CSV/Excel (background)**; download the file from the job form **Staged file or export download** field when state is **Done**.
- Status API: `GET /web/bulk/status?id=<jobId>` (owner only).

## Attachments

Binary uploads use `sys.attachment` and the kernel filestore. Optional virus scanning: register `orm.SetAttachmentScanner` at process startup (see `core/orm/attachment_binary.go`).
