package models

import "sumeru/core/sdk"

// ResConfigSettings is a transient settings shell model extended by addons (not linked from base menus).
type ResConfigSettings struct {
	sdk.Model `sumeru:"model=res.config.settings"`
}
