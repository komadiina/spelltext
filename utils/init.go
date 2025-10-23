package utils

import (
	"time"

	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
)

func InitClientConn(logger *logging.Logger, target string, credentials grpc.DialOption, backoff int, maxRetries int) (*grpc.ClientConn, error) {
	try := 1
	for {
		conn, err := grpc.NewClient(target, credentials)

		if err != nil && try >= maxRetries {
			return nil, err
		} else if err == nil && try < maxRetries {
			return conn, nil
		} else {
			backoff *= 3
			logger.Warnf("unable to connect (%s), backing off... backoff=%ds", err, backoff)
			time.Sleep(time.Duration(backoff) * time.Second)

			try++
		}
	}
}
