package health

import (
	"context"
	"time"

	"github.com/komadiina/spelltext/server/auth/server"
	health "github.com/komadiina/spelltext/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitMonitor(
	s *server.AuthService,
	target string,
	dest health.HealthCheckable,
	onReconnect func(*server.AuthService, *grpc.ClientConn),
) *health.HealthMonitor {
	return &health.HealthMonitor{
		Checker:    dest,
		Logger:     s.Logger,
		Interval:   time.Duration(s.Config.HealthCheckInterval) * time.Second,
		RetryLimit: s.Config.MaxReconnAttempts,
		Target:     target,
		Reconnect: func(ctx context.Context) error { // could reuse InitClientConn somehow, to accept onConnect(), instead of onReconnect()
			try := 1
			backoff := s.Config.Backoff
			for {
				s.Logger.Infof("attempting to reconnect #%d to %s", try, target)
				conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

				if err != nil && try >= s.Config.MaxReconnAttempts {
					return err
				} else if err == nil && conn != nil && try <= s.Config.MaxReconnAttempts {
					s.Logger.Infof("reconnected to %s successfully", target)
					onReconnect(s, conn)
					return nil
				} else {
					s.Logger.Warnf("unable to connect (%s), backing off..., backoff=%ds", err, backoff)

					backoff *= 3
					time.Sleep(time.Duration(backoff) * time.Second)
					try++
				}
			}
		},
	}
}
