package server

import (
	"context"
	"time"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Ping(conn *grpc.ClientConn, ctx context.Context, heartbeatInterval time.Duration) error {
	hc := healthpb.NewHealthClient(conn)

	heartbeatTicker := time.NewTicker(heartbeatInterval)
	for {
		// perform immediate check then wait for ticker
		if err := runHealthCheck(ctx, hc, 3*time.Second); err != nil {
			// connection unhealthy. close and attempt reconnect
			_ = conn.Close()
			m.clearConn()
			heartbeatTicker.Stop()
			break
		}

		select {
		case <-heartbeatTicker.C:
			continue
		case <-ctx.Done():
			heartbeatTicker.Stop()
			_ = conn.Close()
			m.clearConn()
			close(m.stop)
			return
		}
	}

	// small delay before reconnect attempts to avoid hot loop
	time.Sleep(100 * time.Millisecond)
}
