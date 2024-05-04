package pkg

import (
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewSessionManager(pool *pgxpool.Pool) *scs.SessionManager {
	return scs.New()
}
