package server

import (
	"Gin_Event_App_Arc/fabric"
	"Gin_Event_App_Arc/handlers"
	"Gin_Event_App_Arc/middlewares"
	"Gin_Event_App_Arc/repositories"
	"Gin_Event_App_Arc/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (server *Server) Routes() http.Handler {
	g := server.Gin //ye wala to default hai
	//Initiallizing Repos
	userRepo := repositories.InitUserRepo(server.DB)
	eventRepo := repositories.InitEventRepo(server.DB)
	attendeeRepo := repositories.InitAttendeeRepo(server.DB)
	roleRepo := repositories.InitRoleRepo(server.DB)
	userroleRepo := repositories.InitUserRoleRepo(server.DB)

	//Initiallizing Services
	userService := services.NewUserService(userRepo, userroleRepo)
	eventService := services.NewEventService(eventRepo)
	attendeeService := services.NewAttendeeService(attendeeRepo)
	roleService := services.NewRoleService(roleRepo)
	userroleService := services.NewUserRoleService(userroleRepo, roleRepo, userRepo)

	//Hyperledger Network Initialize
	network, gw, err := fabric.Connect()
	if err != nil {
		logrus.Fatalf("❌ Failed to connect to Hyperledger Fabric: %v", err)
	}
	logrus.Println("✅ Connected to Hyperledger Fabric successfully")
	//init userfabricservice
	userfabricService := services.NewUserFabricService(network, gw, "basic")
	//Hyperledger Network Initialize End
	//Initiallizing Handlers
	handlerInitialize := handlers.NewHandlerInstance(userService, eventService, attendeeService, roleService, userroleService, server.RedisClient, server.Hub, userfabricService)
	//Registering Web Socket Endpoint for First Handshake
	g.GET("/ws/users", func(ctx *gin.Context) {
		server.Hub.ServeWS(ctx)
	})
	//Grouping for prefix on API Endpoints

	v1Events := g.Group("/events/v1")
	v1Events.Use(middlewares.JWTAuthMiddleware())
	{
		v1Events.POST("/create", handlerInitialize.CreateEvents) //Because this method takes the context so we do not need to add () with it as we use it in normal methods(Functions tied with specific struct instances)
		v1Events.GET("/getAll/cursor/:cursor", handlerInitialize.GetAllEvents)
		v1Events.GET("/get/:id", handlerInitialize.GetEvents) //here :id is path params
		v1Events.PUT("/update/:id", handlerInitialize.UpdateEvent)
		v1Events.DELETE("/delete/:id", handlerInitialize.DeleteEvent)
		//User Routes
	}
	//Does not need JWT(Public)
	v1Users := g.Group("/users/v1")
	{
		v1Users.POST("/auth/register", handlerInitialize.RegisterUser)
		v1Users.GET("/getAll/page/:page", handlerInitialize.GetAllUsers)
		v1Users.POST("/login", handlerInitialize.LoginUser)

	}

	v1Attendees := g.Group("/attendees/v1")
	v1Attendees.Use(middlewares.JWTAuthMiddleware(), middlewares.RoleAdmin()) //JWT Auth
	{
		v1Attendees.POST("/events/:eventId/attendees/:userId", handlerInitialize.AddAttendeetoEvent)
		v1Attendees.GET("/events/:id/attendees", handlerInitialize.GetAttendeesForEvent)
		v1Attendees.DELETE("/events/:id/attendees/:userId", handlerInitialize.DeleteAttendeeFromEvent)
		v1Attendees.GET("/attendees/:id/events", handlerInitialize.GetEventByAttendee)
		v1Attendees.GET("/allAttendees", handlerInitialize.GetAllAttendeesList)

	}
	v1Roles := g.Group("/roles")
	{
		v1Roles.GET("/get/:id", handlerInitialize.GetRole)
		v1Roles.GET("/getAll", handlerInitialize.GetAllRole)
		v1Roles.POST("/create", handlerInitialize.CreateRole)
		v1Roles.PUT("/update/:id", handlerInitialize.UpdateRole)
		v1Roles.DELETE("/delete/:id", handlerInitialize.DeleteRole)

	}
	v1UserRoles := g.Group("/userRoles")
	{
		v1UserRoles.PUT("/update/userId/:user_id/roleId/:role_id", handlerInitialize.UpdateUserRole)
		v1UserRoles.GET("/getRolesByuserId/:user_id", handlerInitialize.GetRolesByUser)
		v1UserRoles.GET("/getUsersByroleId/:role_id", handlerInitialize.GetUsersByRole)
		v1UserRoles.GET("/getAll", handlerInitialize.GetAll)
	}
	v1UserFabric := g.Group("/userFabric/v1")
	{
		v1UserFabric.GET("/getuser/:id", handlerInitialize.GetUserFabric)
		v1UserFabric.GET("/initledger", handlerInitialize.InitLedger)
		v1UserFabric.PUT("/updateuser/:id", handlerInitialize.UpdateUserFabric)
		v1UserFabric.POST("/adduser", handlerInitialize.InsertUserFabric)
		v1UserFabric.GET("/getfullhistory", handlerInitialize.GetFullHistory)
		v1UserFabric.DELETE("/deleteuser/:id", handlerInitialize.DeleteUserFabric)
		v1UserFabric.GET("/allusers", handlerInitialize.GetAllUsersFabric)
	}
	return g
}
