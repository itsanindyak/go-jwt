package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var googleAuthConfig = &oauth2.Config{
	ClientID: GOOGLE_CLIENT_ID,
	ClientSecret: GOOGLE_CLIENT_SECRET ,
	RedirectURL: GOOGLE_REDIRECT_URL,
	Scopes: []string{"openid", "profile", "email"},
	Endpoint: google.Endpoint,
}

var githubAuthConfig = &oauth2.Config{
	ClientID: GITHUB_CLIENT_ID,
	ClientSecret: GITHUB_CLIENT_SECRET,
	RedirectURL: GITHUB_REDIRECT_URL,
	Scopes: []string{"user:email"},
	Endpoint: github.Endpoint,
}