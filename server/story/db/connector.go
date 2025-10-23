package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/komadiina/spelltext/server/auth/server"
)

func InitializePool(s *server.AuthService, context context.Context, conninfo string, backoff time.Duration, maxRetries int, boFormula func(time.Duration) time.Duration) error {
	try := 1
	for {
		conn, err := pgx.Connect(context, conninfo)

		if err != nil && try >= maxRetries {
			// conn not established, max retries exceeded
			s.Logger.Fatal(err)
		} else if err == nil && try < maxRetries {
			// conn established within maxRetries
			s.Logger.Info("pgpool connection established, creating pool..")
			conn.Close(context)

			pool, err := pgxpool.New(context, fmt.Sprintf(
				"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
				s.Config.PgUser,
				s.Config.PgPass,
				s.Config.PgHost,
				s.Config.PgPort,
				s.Config.PgDbName,
				s.Config.PgSSLMode,
			))

			if err != nil {
				s.Logger.Fatal("unable to create pool", "reason", err)
			} else {
				s.Logger.Info("pgxpool (dpool, via pgpool-ii) initialized")
			}

			s.DbPool = pool

			return nil
		} else {
			// conn not established, backoff
			s.Logger.Warn("failed to establish database connection, backing off...", "reason", err, "backoff_seconds", backoff.Seconds())
			time.Sleep(backoff)
			backoff = boFormula(backoff)
			try++
		}
	}
}
