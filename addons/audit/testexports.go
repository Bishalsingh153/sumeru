package audit

// AuditExportRowForTest is one exported audit row for tests.
type AuditExportRowForTest = auditExportRow

// RetentionDaysFromRecordForTest exposes audit cleanup preset parsing for tests.
func RetentionDaysFromRecordForTest(rec map[string]interface{}) (int, error) {
	return retentionDaysFromRecord(rec)
}
