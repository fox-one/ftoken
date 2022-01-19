package session

import (
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/mixin-sdk-go"
	"golang.org/x/sync/singleflight"
)

func New(userz core.UserService, capacity int, issuers []string) core.Session {
	var s core.Session = &session{
		userz:   userz,
		issuers: issuers,
		sf:      &singleflight.Group{},
	}

	if capacity > 0 {
		return Cache(s, capacity)
	}

	return s
}

type session struct {
	userz   core.UserService
	sf      *singleflight.Group
	issuers []string
}

func (s *session) Login(ctx context.Context, accessToken string) (*core.User, error) {
	user, err, _ := s.sf.Do(accessToken, func() (interface{}, error) {
		var claim struct {
			jwt.StandardClaims
			Scope string `json:"scp,omitempty"`
		}
		_, _ = jwt.ParseWithClaims(accessToken, &claim, nil)

		if jti := claim.Id; govalidator.IsUUID(jti) {
			ctx = mixin.WithRequestID(ctx, jti)
		}

		user, err := s.userz.Login(ctx, accessToken)
		if err != nil {
			return nil, err
		}

		return user, nil
	})

	if err != nil {
		return nil, err
	}

	return user.(*core.User), nil
}
