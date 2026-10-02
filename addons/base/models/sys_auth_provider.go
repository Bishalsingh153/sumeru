package models

import (
	"sumeru/core/sdk"
)

type SysAuthProvider struct {
	sdk.Model `sumeru:"model=sys.auth.provider"`

	Name         sdk.String  `sumeru:"required,string=Name"`
	ProviderType sdk.String  `sumeru:"required,string=Type,default=oidc,selection=oidc:OpenID Connect,oauth2:OAuth2"`
	ClientID     sdk.String  `sumeru:"required,string=Client ID"`
	ClientSecret sdk.String  `sumeru:"string=Client Secret"`
	IssuerURL    sdk.String  `sumeru:"string=Issuer URL"`
	AuthorizeURL sdk.String  `sumeru:"string=Authorize URL"`
	TokenURL     sdk.String  `sumeru:"string=Token URL"`
	JwksURL      sdk.String  `sumeru:"string=JWKS URL"`
	Scopes       sdk.String  `sumeru:"string=Scopes,default=openid email profile"`
	Enabled      sdk.Boolean `sumeru:"string=Enabled,default=false"`
	ButtonLabel  sdk.String  `sumeru:"string=Button Label"`
	LinkPolicy   sdk.String  `sumeru:"string=Unknown Users,default=deny,selection=deny:Deny,link:Link by email,jit_internal:JIT internal,jit_portal:JIT portal"`
}
