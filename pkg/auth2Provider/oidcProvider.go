package auth2provider

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/itsanindyak/go-jwt/config"
	"golang.org/x/oauth2"
)

type OIDCProvider struct {
	name     string
	issuer   string
	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier
	provider *oidc.Provider
}

func NewOIDCProvider(name, issuer,clientID,clientSecret,redirectURL string,scopes []string) *OIDCProvider {
	ctx := context.Background()

	p, err := oidc.NewProvider(ctx, issuer)

	if err != nil {
		panic(fmt.Sprintf("[%s] failed to init OIDC provider (%s): %v", name, issuer, err.Error()))
	}

	if len(scopes) ==0 {
		scopes = []string{oidc.ScopeOpenID,"profile","email"}
	}

	config := &oauth2.Config{
		ClientID: clientID,
		ClientSecret: clientSecret,
		Endpoint: p.Endpoint(),
		Scopes: scopes,
		RedirectURL: redirectURL,
	}

	return &OIDCProvider{
		name: name,
		issuer: issuer,
		config: config,
		verifier: p.Verifier(&oidc.Config{ClientID: clientID}),
		provider: p,
	}
}

func (o *OIDCProvider) Name() string{
	return o.name
}

func (o *OIDCProvider) AuthURL(state string) string{
	return o.config.AuthCodeURL(state)
}

// func (o *OIDCProvider) HandleCallback()(*User,error){

// }

func init(){
	if config.GOOGLE_CLIENT_ID != ""{
		Register("google",func() Provider {
			return NewOIDCProvider(
				"google",
				config.GOOGLE_OIDC_ISSUER,
				config.GOOGLE_CLIENT_ID,
				config.GOOGLE_CLIENT_SECRET,
				config.GOOGLE_REDIRECT_URL,
				nil,

			)
		})
	}
}