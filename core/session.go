package core

import (
	"context"
)

type (
	User struct {
		MixinID string `json:"mixin_id,omitempty"`
		Role    string `json:"role,omitempty"`
		Lang    string `json:"lang,omitempty"`
		Name    string `json:"name,omitempty"`
		Avatar  string `json:"avatar,omitempty"`
	}

	UserService interface {
		Find(ctx context.Context, mixinID string) (*User, error)
		Login(ctx context.Context, token string) (*User, error)
	}

	Session interface {
		// Login return user mixin id
		Login(ctx context.Context, accessToken string) (*User, error)
	}
)
