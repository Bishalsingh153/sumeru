package models

import (
	"sumeru/core/sdk"
)

type CoreUserMail struct {
	sdk.Model `sumeru:"inherit=core.user"`

	NotifyEmail sdk.Boolean `sumeru:"string=Email Notifications,default=true,column=notify_email"`
}
