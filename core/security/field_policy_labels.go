package security

import "strings"

// FieldKernelPolicyTags returns read-only matrix legend tags for kernel field policy.
func FieldKernelPolicyTags(model, field string) []string {
	model = strings.TrimSpace(model)
	field = strings.TrimSpace(field)
	if model == "" || field == "" {
		return nil
	}
	var tags []string
	if ReadRedactFields(model)[field] {
		tags = append(tags, "Kernel read redact")
	}
	if WriteDenyDirectFields(model)[field] {
		tags = append(tags, "Kernel write deny")
	}
	if WriteDenyUnlessSysFields(model)[field] {
		tags = append(tags, "Kernel sysadmin write")
	}
	return tags
}
