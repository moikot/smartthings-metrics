package main

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

type AuthInfoWriter struct {
	token string
}

func NewAuthInfoWriter(token string) runtime.ClientAuthInfoWriter {
	return &AuthInfoWriter{
		token: token,
	}
}

func (a *AuthInfoWriter) AuthenticateRequest(r runtime.ClientRequest, reg strfmt.Registry) error {
	return r.SetHeaderParam("Authorization", "Bearer "+a.token)
}
