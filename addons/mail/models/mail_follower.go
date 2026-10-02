package models

import (
	"sumeru/core/sdk"
)

type MailFollower struct {
	sdk.Model `sumeru:"model=mail.follower"`

	ResModel    sdk.String                    `sumeru:"required,index,column=res_model,string=Document Model"`
	ResID       sdk.Integer                   `sumeru:"required,index,column=res_id,string=Document ID"`
	UserID      sdk.Many2One[CoreUser]        `sumeru:"required,index,string=User"`
	PartnerID   sdk.Many2One[CorePartner]     `sumeru:"index,string=Partner"`
	SubtypeIds  sdk.Many2Many[MailMessageSubtype] `sumeru:"string=Subtypes,table=mail_follower_subtype_rel,left=follower_id,right=subtype_id"`
	Active      sdk.Boolean                   `sumeru:"string=Active,default=true"`
}
