package pkg

import (
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

const SESSION_KEY_MESSAGE = "message"
const LOGGED_IN_USERNAME = "logged_in_username"

func NewSessionManager(pool *pgxpool.Pool) *scs.SessionManager {
	return scs.New()
}
