package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	*http.Server
	logger *logrus.Logger
}

func NewServer(gin *gin.Engine, logger *logrus.Logger) *Server {
	// 在日志中输出挂载的路由
	for _, v := range gin.Routes() {
		logger.Debugf("controller %s mounted at %s", v.Method, v.Path)
	}

	server := &Server{
		Server: &http.Server{
			Addr:    net.JoinHostPort("0.0.0.0", "8123"),
			Handler: gin,
		},
		logger: logger,
	}

	return server
}

func (s *Server) Run(ctx context.Context, group *errgroup.Group, shutdownHttpServer <-chan struct{}) {
	group.Go(func() error {
		s.logger.Infof("listening and serving HTTP on %s", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			return err
		}

		return nil
	})

	group.Go(func() error {
		select {
		case <-ctx.Done():
		case <-shutdownHttpServer:
		}

		s.logger.Info("gracefully shutting down service...")
		closeCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.Shutdown(closeCtx); err != nil && err != http.ErrServerClosed {
			s.logger.Errorf("shutdown server failed: %v", err)
			return err
		}

		return nil
	})
}
