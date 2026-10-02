package models

import (
	"sumeru/core/sdk"
)

type MailNotification struct {
	sdk.Model `sumeru:"model=mail.notification"`

	UserID       sdk.Many2One[CoreUser]    `sumeru:"required,index,string=User"`
	MessageID    sdk.Many2One[MailMessage] `sumeru:"required,index,string=Message"`
	IsRead       sdk.Boolean               `sumeru:"string=Read,default=false"`
	ReadDate     sdk.DateTime              `sumeru:"string=Read Date"`
	RecordModel  sdk.String                `sumeru:"index,column=record_model,string=Record Model"`
	RecordID     sdk.Integer               `sumeru:"index,column=record_id,string=Record ID"`
	CompanyID    sdk.Many2One[CoreCompany] `sumeru:"index,string=Company"`
	CreateDate   sdk.DateTime              `sumeru:"required,string=Created"`
}
