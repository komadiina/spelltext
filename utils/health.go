package utils

import (
	"context"
	"time"

	pbHealth "github.com/komadiina/spelltext/proto/health"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
)

type HealthCheckable interface {
	Check(ctx context.Context, req *pbHealth.HealthCheckRequest, opts ...grpc.CallOption) (*pbHealth.HealthCheckResponse, error)
}

type ReconnectFunc func(ctx context.Context) error

type HealthMonitor struct {
	Checker    HealthCheckable
	Logger     *logging.Logger
	Interval   time.Duration
	RetryLimit int
	Reconnect  ReconnectFunc
	Target     string
}

func (m *HealthMonitor) Run(ctx context.Context) {
	if m.Interval <= 0 {
		m.Interval = 10 * time.Second
	}

	retries := m.RetryLimit

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(m.Interval):
			resp, err := m.Checker.Check(ctx, &pbHealth.HealthCheckRequest{})

			if err != nil {
				m.Logger.Error("server timeout", "reason", err)

				if retries <= 0 {
					m.Logger.Warnf("service %s unhealthy, attempting reconnect", m.Target)
					if m.Reconnect == nil {
						m.Logger.Error("no reconnect function provided")
						return
					}
					if err := m.Reconnect(ctx); err != nil {
						m.Logger.Errorf("failed to reconnect, service %s down. reason=%v", m.Target, err)
						return
					}

					// reset retries
					retries = m.RetryLimit
				} else {
					retries--
				}
			} else {
				// healthy, so reset retries
				retries = m.RetryLimit
				m.Logger.Info("server is back up, healthy. service=%s, status=%v", m.Target, resp.Status)
			}
		}
	}
}
