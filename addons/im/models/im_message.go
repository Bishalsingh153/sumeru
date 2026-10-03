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
	IsRead      sdk.Boolean            `sumeru:"default=false,index,column=is_read,string=Read"`
	ResModel    sdk.String             `sumeru:"index,column=res_model,string=Related document model"`
	ResID       sdk.Integer            `sumeru:"index,string=Related document id"`
}
