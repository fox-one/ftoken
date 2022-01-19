package session

import (
	"context"

	"github.com/bluele/gcache"
	"github.com/fox-one/ftoken/core"
)

type cacheSession struct {
	core.Session
	tokens gcache.Cache
}

func Cache(session core.Session, capacity int) core.Session {
	return &cacheSession{
		Session: session,
		tokens:  gcache.New(capacity).LRU().Build(),
	}
}

func (s *cacheSession) Login(ctx context.Context, accessToken string) (*core.User, error) {
	if user, err := s.tokens.Get(accessToken); err == nil {
		return user.(*core.User), nil
	}

	user, err := s.Session.Login(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	_ = s.tokens.Set(accessToken, user)
	return user, nil
}
