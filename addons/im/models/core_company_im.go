package models

import (
	"sumeru/core/sdk"
)

type CoreCompanyIM struct {
	sdk.Model `sumeru:"inherit=core.company"`

	ImEnabled sdk.Boolean `sumeru:"string=Internal chat,default=true,column=im_enabled"`
}
