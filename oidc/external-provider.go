package oidc

import (
	"context"
	"errors"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var ExternalProvider externalProvider
var OidcProvider *oidc.Provider
var TokenVerifier *oidc.IDTokenVerifier

type externalProvider struct {
	oauth2.Config
}

func (p *externalProvider) Configure(issuer string, clientID string, clientSecret string, redirectURI string) error {
	if strings.TrimSpace(issuer) == "" {
		return errors.New("empty issuer")
	}

	if strings.TrimSpace(clientID) == "" {
		return errors.New("empty clientID")
	}

	if strings.TrimSpace(clientSecret) == "" {
		return errors.New("empty clientSecret")
	}

	if strings.TrimSpace(redirectURI) == "" {
		return errors.New("empty redirectURI")
	}

	p.ClientID = clientID
	p.ClientSecret = clientSecret
	p.RedirectURL = redirectURI
	p.Scopes = []string{oidc.ScopeOpenID, "profile", "email"}

	var err error
	OidcProvider, err = oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		return err
	}

	TokenVerifier = OidcProvider.Verifier(&oidc.Config{ClientID: clientID})

	p.Endpoint = OidcProvider.Endpoint()
	return nil
}
