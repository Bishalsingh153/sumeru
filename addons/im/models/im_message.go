package models

import (
	"sumeru/core/sdk"
)

type ImMessage struct {
	sdk.Model `sumeru:"model=im.message"`

	SenderID    sdk.Many2One[CoreUser] `sumeru:"required,index,string=From"`
	RecipientID sdk.Many2One[CoreUser] `sumeru:"required,index,string=To"`
	Body        sdk.Text               `sumeru:"string=Message"`
	CreateDate  sdk.DateTime           `sumeru:"required,index,string=Sent"`
}
