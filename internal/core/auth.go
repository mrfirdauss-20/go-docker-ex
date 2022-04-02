package core

import (
	"context"

	"gopkg.in/validator.v2"
)

type Auth interface {
	ValidateAPIKey(ctx context.Context, apiKey string) error
}

type AuthConfig struct {
	APIKey string `validate:"nonzero"`
}

func (c AuthConfig) Validate() error {
	return validator.Validate(c)
}

func NewAuth(cfg AuthConfig) (Auth, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &auth{apiKey: cfg.APIKey}, nil
}

type auth struct {
	apiKey string
}

func (a *auth) ValidateAPIKey(ctx context.Context, apiKey string) error {
	if a.apiKey != apiKey {
		return ErrInvalidAPIKey
	}
	return nil
}
