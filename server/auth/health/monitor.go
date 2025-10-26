package health

import (
	"context"
	"errors"
	"time"

	"github.com/komadiina/spelltext/server/auth/server"
	health "github.com/komadiina/spelltext/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func InitMonitor(
	s *server.AuthService,
	target string,
	dest health.HealthCheckable,
	onReconnect func(*server.AuthService, *grpc.ClientConn),
) *health.HealthMonitor {
	s.Logger.Infof("initializing health monitor for %s", target)

	return &health.HealthMonitor{
		Checker:    dest,
		Logger:     s.Logger,
		Interval:   time.Duration(s.Config.HealthCheckInterval) * time.Second,
		RetryLimit: s.Config.MaxReconnAttempts,
		Target:     target,
		Reconnect: func(ctx context.Context) error { // could reuse InitClientConn somehow, to accept onConnect(), instead of onReconnect()
			backoff := s.Config.Backoff

			bo := func() {
				time.Sleep(time.Duration(backoff) * time.Second)
				backoff *= 2
			}

			for {
				select {
				case <-ctx.Done():
					return errors.New("context canceled")
				default:
					conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						s.Logger.Warnf("dial error: %v", err)
						bo()
					}

					if conn.GetState() != connectivity.Ready {
						s.Logger.Infof("attempting to reconnect to %s", target)
						conn.Connect()
					}

					// wait for statechange
					time.Sleep(500 * time.Millisecond)
					if conn.GetState() == connectivity.Ready {
						onReconnect(s, conn)
						return nil
					}

					bo()
				}
			}
		},
	}
}
