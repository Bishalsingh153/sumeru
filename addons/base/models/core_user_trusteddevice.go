package models

import (
	"sumeru/core/sdk"
)

type CoreUserTrustedDevice struct {
	sdk.Model `sumeru:"model=core.user.trusteddevice"`

	UserID    sdk.Many2One[CoreUser] `sumeru:"required,index,string=User"`
	TokenHash sdk.String             `sumeru:"required,index,string=Token Hash,column=token_hash"`
	UserAgent sdk.String             `sumeru:"string=User Agent,column=user_agent"`
	ExpiresAt sdk.DateTime           `sumeru:"required,index,string=Expires At,column=expires_at"`
}
