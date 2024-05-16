package pkg

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/net/context"
)

const SESSION_KEY_MESSAGE = "message"
const LOGGED_IN_USERNAME = "logged_in_username"

type SessionManager interface {
	Exists(ctx context.Context, key string) bool
	Put(ctx context.Context, key string, val interface{})
	Remove(ctx context.Context, key string)
	PopString(ctx context.Context, key string) string
	GetString(ctx context.Context, key string) string
	LoadAndSave(next http.Handler) http.Handler
	RenewToken(ctx context.Context) error
}

type SessionManagerImpl struct {
	scs *scs.SessionManager
}

func NewSessionManagerImpl() SessionManagerImpl {
	return SessionManagerImpl{
		scs: scs.New(),
	}
}

func (s SessionManagerImpl) Exists(ctx context.Context, key string) bool {
	return s.scs.Exists(ctx, key)
}

func (s SessionManagerImpl) PopString(ctx context.Context, key string) string {
	return s.scs.PopString(ctx, key)
}

func (s SessionManagerImpl) GetString(ctx context.Context, key string) string {
	return s.scs.GetString(ctx, key)
}

func (s SessionManagerImpl) LoadAndSave(next http.Handler) http.Handler {
	return s.scs.LoadAndSave(next)
}

func (s SessionManagerImpl) RenewToken(ctx context.Context) error {
	return s.scs.RenewToken(ctx)
}

func (s SessionManagerImpl) Put(ctx context.Context, key string, val interface{}) {
	s.scs.Put(ctx, key, val)
}

func (s SessionManagerImpl) Remove(ctx context.Context, key string) {
	s.scs.Remove(ctx, key)
}
