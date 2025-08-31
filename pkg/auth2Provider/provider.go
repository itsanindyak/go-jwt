package auth2provider

import "net/http"

type User struct {
	ID       string
	Email    string
	Name     string
	Picture  string
	Provider string
}

type Provider interface {
	Name() string
	AuthURL(state string) string
	// HandleCallback() (*User, error)
}