package models

import (
	"sumeru/core/sdk"
)

type MailMessageSubtype struct {
	sdk.Model `sumeru:"model=mail.message.subtype"`

	Name        sdk.String  `sumeru:"required,unique,string=Technical Name"`
	Description sdk.String  `sumeru:"string=Description"`
	Internal    sdk.Boolean `sumeru:"string=Internal,default=false"`
	Default     sdk.Boolean `sumeru:"string=Default,default=false"`
	Hidden      sdk.Boolean `sumeru:"string=Hidden,default=false"`
	Sequence    sdk.Integer `sumeru:"string=Sequence,default=10"`
}
