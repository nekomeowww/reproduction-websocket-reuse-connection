package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	gin := NewGin(logger)
	gin.Handle(http.MethodGet, "ws/v1/chat/:roomId", HandleWebsocket(logger))

	{
		server := NewServer(gin, logger)

		shutdownHTTPServerChannel := make(chan struct{})
		group, ctx := errgroup.WithContext(context.Background())

		group.Go(func() error {
			osCh := make(chan os.Signal, 1)
			signal.Notify(osCh, os.Interrupt, syscall.SIGTERM)
			select {
			case <-ctx.Done():
			case <-osCh:
			}

			shutdownHTTPServerChannel <- struct{}{}
			return nil
		})

		server.Run(ctx, group, shutdownHTTPServerChannel)
		err := group.Wait()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}
}
