package swcmeta

import "sumeru/core/orm"

type UserSecurityGroupRow struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Selected   bool   `json:"selected"`
	IsUserType bool   `json:"isUserType"`
}

type UserSecurityCompanyRow struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Selected bool   `json:"selected"`
}

type UserSecurityPayload struct {
	CanEdit       bool                     `json:"canEdit"`
	IsNew         bool                     `json:"isNew"`
	Groups        []UserSecurityGroupRow   `json:"groups"`
	Companies     []UserSecurityCompanyRow `json:"companies"`
	PasswordLabel string                   `json:"passwordLabel"`
}

func UserSecurityFromORM(m orm.UserSecurityMeta) UserSecurityPayload {
	out := UserSecurityPayload{
		CanEdit:       m.CanEdit,
		IsNew:         m.IsNew,
		PasswordLabel: m.PasswordLabel,
	}
	for _, g := range m.Groups {
		out.Groups = append(out.Groups, UserSecurityGroupRow{
			ID: g.ID, Name: g.Name, Selected: g.Selected, IsUserType: g.IsUserType,
		})
	}
	for _, c := range m.Companies {
		out.Companies = append(out.Companies, UserSecurityCompanyRow{
			ID: c.ID, Name: c.Name, Selected: c.Selected,
		})
	}
	return out
}
