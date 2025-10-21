package main

import (
	"Gin_Event_App_Arc/config"
	"Gin_Event_App_Arc/db"
	"Gin_Event_App_Arc/kafka"
	"Gin_Event_App_Arc/repositories"
	"Gin_Event_App_Arc/server"
	"Gin_Event_App_Arc/services"
	userpb "Gin_Event_App_Arc/src/proto"
	"Gin_Event_App_Arc/ws"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	_ "modernc.org/sqlite"
)

func main() {
	//New Code for Starting Server
	newConfig := config.NewConfig()
	redisClient := config.NewRedisClient()
	database := db.InitDB(newConfig.DB)
	httpsV := gin.Default()
	srv := &http.Server{
		Handler: httpsV,
	}
	//Creating and Storing Sessions
	store := cookie.NewStore([]byte("SecretKey-For-Auth"))
	store.Options(sessions.Options{
		MaxAge:   60 * 30, // 30 minutes (in seconds)
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true when using HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	httpsV.Use(sessions.Sessions("NewSession", store))
	hub := ws.NewHub()
	go hub.Run()
	app := server.NewServer(newConfig, redisClient, httpsV, database, hub)
	app.Routes()
	// Not MINE
	// Running the server on another Goroutine
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- app.Run(newConfig.HTTP.Port)
	}()

	//Kafka Init
	kafka.InitProducer("localhost:9092")
	defer kafka.CloseProducer(context.Background())
	//Kafka End
	//gRPC Start Code
	userRepo := repositories.InitUserRepo(database)
	userRoleRepo := repositories.InitUserRoleRepo(database)
	userService := &services.UserService{
		UserServ:     userRepo,
		UserRoleServ: userRoleRepo,
	}
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &services.UserProto{Serviceprot: userService})
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logrus.Fatalf("Failed to Listen %v", err.Error())
	}
	go func() {
		logrus.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			logrus.Fatalf("Failed to Serve gRPC : %v", err.Error())
		}
	}()
	//End gRPC Code
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Println("Server is shutting down...")
	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Attempt to gracefully shutdown the server
	if err := app.Shutdown(ctx, srv, grpcServer); err != nil {
		logrus.Printf("Failed to gracefully shutdown the server: %v\n", err)
	}
	// // Wait for the server to shutdown
	select {
	case err := <-serverErr:
		if err != nil {
			logrus.Printf("Server error: %v\n", err)
		}
	case <-ctx.Done():
		logrus.Println("Graceful shutdown timeout exceeded")
	}

	logrus.Println("Server stopped.")

	//Old Code for Starting Server

	// dsn := "host=localhost user=postgres password=Harshal@2003 dbname=eventsapp port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	// // db, err := gorm.Open(sqlite.Open("./data.db"), &gorm.Config{})
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal("Cannot Connect to DB", err)
	// }

	// err2 := db.AutoMigrate(
	// 	&models.Event{},
	// 	&models.User{},
	// 	&models.Attendees{},
	// )
	// if err2 != nil {
	// 	log.Fatal("Failed to migrate tables:", err2)
	// }
	// model := models.NewModels(db)

	// redisClient := config.NewRedisClient()
	// app := &models.Application{
	// 	Port:   env.GetENVIntegers("PORT", 8080),
	// 	Models: model,
	// 	Cache:  redisClient,
	// }
	// var tempRout *routes.RouteStruct

	// server := &http.Server{
	// 	Addr:         fmt.Sprintf(":%d", app.Port),
	// 	Handler:      tempRout.Routes(),
	// 	IdleTimeout:  time.Minute,
	// 	ReadTimeout:  10 * time.Second,
	// 	WriteTimeout: 30 * time.Second,
	// }

	// if err := server.ListenAndServe(); err != nil {
	// 	log.Fatal("Error for Server", err)
	// }

	// if err := app.Serve(); err != nil {
	// 	log.Fatal("Hello", err)
	// }
}
