package server

import (
	"Gin_Event_App_Arc/config"
	"Gin_Event_App_Arc/ws"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	Cfg         config.Config
	Gin         *gin.Engine
	DB          *gorm.DB
	RedisClient *config.RedisClient
	Hub         *ws.Hub
}

func NewServer(cfg config.Config, redis *config.RedisClient, httpHandle *gin.Engine, db *gorm.DB, hub *ws.Hub) *Server {

	return &Server{
		Cfg:         cfg,
		Gin:         httpHandle,
		DB:          db,
		RedisClient: redis,
		Hub:         hub,
	}
}

func (server Server) Run(addr string) error {
	return server.Gin.Run(":" + addr)
}

func (s Server) Shutdown(ctx context.Context, srv *http.Server, grpcServer *grpc.Server) error {

	// Gracefully shutdown the session of gin
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("error closing session of gin: %v", err)
	}

	// Gracefully close the Kafka client
	// if err := s.KafkaClient.Close(); err != nil {
	// 	return fmt.Errorf("error closing Kafka client: %v", err)
	// }

	// Gracefully close the redis client
	if err := s.RedisClient.Close(); err != nil {
		return fmt.Errorf("error closing redis client: %v", err)
	}

	db, err := s.DB.DB()
	if err != nil {
		return err
	}

	// Attempt to gracefully close the database connection
	if err := db.Close(); err != nil {
		return err
	}
	grpcServer.GracefulStop()

	return nil
}
