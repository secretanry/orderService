package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"wb-L0/modules/config"
	"wb-L0/routing"
)

type Server struct {
	Serv *http.Server
}

func (s *Server) Init(errChan chan error) error {
	switch config.GetConfig().RunMode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		return fmt.Errorf("run mode %s not supported", config.GetConfig().RunMode)
	}
	r := gin.Default()
	routing.MountSystemRoutes(r)
	routing.MountPurchasesRoutes(r)
	routing.MountFrontRoutes(r)
	addr := fmt.Sprintf("0.0.0.0:%d", config.GetConfig().AppPort)
	s.Serv = &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		if err := s.Serv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()
	return nil
}

func (s *Server) SuccessfulMessage() string {
	return fmt.Sprintf("Server successfully started on address %s", s.Serv.Addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Serv.Shutdown(ctx)
}
