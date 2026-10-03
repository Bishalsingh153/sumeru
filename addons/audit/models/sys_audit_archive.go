package models

import (
	"sumeru/core/sdk"
)

type SysAuditArchive struct {
	sdk.Model `sumeru:"model=sys.audit.archive"`

	Name         sdk.String                  `sumeru:"required,string=Name"`
	State        sdk.String                  `sumeru:"required,index,string=State,default=stored,selection=stored:Stored,failed:Failed,purged:Purged"`
	PeriodFrom   sdk.DateTime                `sumeru:"string=Period from"`
	PeriodTo     sdk.DateTime                `sumeru:"string=Period to"`
	RowCount     sdk.Integer                 `sumeru:"string=Rows"`
	ByteSize     sdk.Integer                 `sumeru:"string=Bytes"`
	SHA256       sdk.String                  `sumeru:"string=SHA-256"`
	AttachmentID sdk.Many2One[SysAttachment] `sumeru:"string=Download"`
	Notes        sdk.Text                    `sumeru:"string=Notes"`
	ErrorMessage sdk.Text                    `sumeru:"string=Error"`
	CreateDate   sdk.DateTime                `sumeru:"string=Created"`
}
