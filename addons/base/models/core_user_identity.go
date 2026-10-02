package models

import (
	"sumeru/core/sdk"
)

type CoreUserIdentity struct {
	sdk.Model `sumeru:"model=core.user.identity"`

	ProviderID sdk.Many2One[SysAuthProvider] `sumeru:"required,index,string=Provider"`
	Subject    sdk.String                    `sumeru:"required,index,string=Subject"`
	UserID     sdk.Many2One[CoreUser]        `sumeru:"required,index,string=User"`
}
