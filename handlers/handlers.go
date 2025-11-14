package handlers

import (
	"Gin_Event_App_Arc/config"
	"Gin_Event_App_Arc/models"
	"Gin_Event_App_Arc/response"
	"Gin_Event_App_Arc/services"
	"Gin_Event_App_Arc/utils"
	"Gin_Event_App_Arc/ws"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type HandlerStruct struct {
	UserHandler       services.UserServiceI
	EventHandler      services.EventServiceI
	AttendeeHandler   services.AttendeeServiceI
	RoleHandler       services.RoleServiceI
	UserRoleHandler   services.UserRoleServiceI
	Cache             *config.RedisClient
	Hub               *ws.Hub
	UserFabricHandler services.UserFabricServiceI
	FileUploadHandler services.FileServiceI
}

func NewHandlerInstance(userHandler services.UserServiceI, eventHandler services.EventServiceI, attendeeHandler services.AttendeeServiceI, roleHandler services.RoleServiceI, userroleHandler services.UserRoleServiceI, cache *config.RedisClient, hub *ws.Hub, userFabricHandler services.UserFabricServiceI, fileUploadHandler services.FileServiceI) *HandlerStruct {
	return &HandlerStruct{
		UserHandler:       userHandler,
		EventHandler:      eventHandler,
		AttendeeHandler:   attendeeHandler,
		RoleHandler:       roleHandler,
		UserRoleHandler:   userroleHandler,
		Cache:             cache,
		Hub:               hub,
		UserFabricHandler: userFabricHandler,
		FileUploadHandler: fileUploadHandler,
	}
}

//CONTROLLERS

// Create Event Controller
func (handler *HandlerStruct) CreateEvents(c *gin.Context) {
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := handler.EventHandler.Insert(&event) //Calling to the service

	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Response(c, http.StatusCreated, event)
}

// Get Event By ID Controller
func (handler *HandlerStruct) GetEvents(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//Caching check
	ctx := context.Background()
	key := fmt.Sprintf("eventByID:%d", id) // -> "event:1"
	cached, errcach := handler.Cache.Client.Get(ctx, key).Result()
	if errcach != redis.Nil {
		var event *models.Event
		json.Unmarshal([]byte(cached), &event)
		response.Response(c, http.StatusOK, event)
		return
	}

	eventFromDB, err1 := handler.EventHandler.Get(id)

	if eventFromDB == nil {
		response.ErrorResponse(c, http.StatusNotFound, "No Such events exists")
		return
	}

	if reflect.DeepEqual(eventFromDB, models.Event{}) {
		response.ErrorResponse(c, http.StatusNotFound, "No Such events exists")
	}

	if err1 != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err1.Error())
	}

	data, _ := json.Marshal(eventFromDB)
	handler.Cache.Client.Set(ctx, key, data, 2*time.Minute)
	response.SuccessResponse(c, eventFromDB)

}

// GetAll Controller
func (handler *HandlerStruct) GetAllEvents(c *gin.Context) {
	cursor, errp := strconv.Atoi(c.Param("cursor"))
	if errp != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid Page No")
		return
	}
	allEvents, nextCursor, err := handler.EventHandler.GetAll(cursor)
	if len(allEvents) == 0 {
		response.ErrorResponse(c, http.StatusNotFound, "Page Not Found")
		return
	}
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Error Fetching List of Events")
		return
	}
	c.JSON(http.StatusOK, gin.H{"next_cursor": nextCursor, "Result": allEvents})

}

// Update Controller
func (handler *HandlerStruct) UpdateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid Event ID in request Params")
		return
	}
	existingEvent, err := handler.EventHandler.Get(id)

	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	if existingEvent == nil {
		response.ErrorResponse(c, http.StatusNotFound, "No such Event Exists ")
		return
	}

	updatedEvent := &models.Event{} //Ek nya Struct ka instance create kiya

	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	updatedEvent.Id = id
	updatedEvent, err2 := handler.EventHandler.Update(updatedEvent, id)

	if err2 != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err2.Error())
		return
	}
	response.SuccessResponse(c, updatedEvent)
}

func (handler *HandlerStruct) DeleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid Event ID in request Params")
		return
	}
	if err := handler.EventHandler.Delete(id); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	response.Response(c, http.StatusNoContent, nil)

}

func (handler *HandlerStruct) AddAttendeetoEvent(c *gin.Context) {
	//UserId valid
	eventId, err := strconv.Atoi(c.Param("eventId"))

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid Event ID")
		return
	}
	//Id Valid
	id, err1 := strconv.Atoi(c.Param("userId"))

	if err1 != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid User ID")
		return
	}
	//Check if event Exists or not
	event, err := handler.EventHandler.Get(eventId)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Error fetching event")
		return
	}
	if event == nil {
		response.ErrorResponse(c, http.StatusNotFound, "Event Not Present")
		return
	}
	//Check if user Exists or not
	userNow, err2 := handler.EventHandler.Get(id)
	if err2 != nil {
		response.ErrorResponse(c, http.StatusNotFound, err2.Error())
		return
	}
	if userNow == nil {
		response.ErrorResponse(c, http.StatusNotFound, "No User Found")
		return
	}
	//Check if the entry is not duplicate
	attendeeCheck, err2 := handler.AttendeeHandler.GetByEventAndAttendee(event.Id, userNow.Id)
	if err2 != nil {
		response.ErrorResponse(c, http.StatusNotFound, err2.Error())
		return
	}
	if attendeeCheck != nil {
		response.ErrorResponse(c, http.StatusConflict, "Already Part of this event")
		return
	}
	//Now we can insert the Event-User in Attendees Table
	attendee := models.Attendees{
		EventId: eventId,
		UserId:  id,
	}

	_, err = handler.AttendeeHandler.Insert(&attendee)

	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Response(c, http.StatusCreated, nil)

}

func (handler *HandlerStruct) GetAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid ID params in URL")
		return
	}

	//Caching check
	ctx := context.Background()
	key := fmt.Sprintf("attendeeByEvent:%d", id) // -> "event:1"
	cached, errcach := handler.Cache.Client.Get(ctx, key).Result()
	if errcach != redis.Nil {
		var event []*models.User
		json.Unmarshal([]byte(cached), &event)
		c.JSON(http.StatusOK, gin.H{"data": event, "cached": true})
		return
	}

	attendeesByEvent, err := handler.AttendeeHandler.GetAttendeesByEvent(id)

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	data, _ := json.Marshal(attendeesByEvent)
	handler.Cache.Client.Set(ctx, key, data, 2*time.Minute)
	response.SuccessResponse(c, attendeesByEvent)

}

func (handler *HandlerStruct) DeleteAttendeeFromEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid ID param in URL")
		return
	}

	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid User ID in params")
		return
	}

	errNew := handler.AttendeeHandler.Delete(eventId, userID)

	if errNew != nil {
		response.ErrorResponse(c, http.StatusBadRequest, errNew.Error())
	}
	response.Response(c, http.StatusNoContent, nil)
}

func (handler *HandlerStruct) GetEventByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid ID param in URL")
		return
	}

	ctx := context.Background()
	key := fmt.Sprintf("eventByAttendee:%d", id) // -> "event:1"
	cached, errcach := handler.Cache.Client.Get(ctx, key).Result()
	if errcach != redis.Nil {
		var event []*models.Event
		json.Unmarshal([]byte(cached), &event)
		c.JSON(http.StatusOK, gin.H{"data": event, "cached": true})
		return
	}
	event, err := handler.EventHandler.GetEventByAttendee(id)

	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	data, _ := json.Marshal(event)
	handler.Cache.Client.Set(ctx, key, data, 2*time.Minute)
	response.SuccessResponse(c, event)
}

func (handler *HandlerStruct) GetAllUsers(c *gin.Context) {
	fmt.Println("Checking if Status hai", http.StatusOK)
	pageNo, errp := strconv.Atoi(c.Param("page"))
	if errp != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid Page No in the URL")
		return
	}

	//Caching check
	ctx := context.Background()
	key := fmt.Sprintf("pageResult:%d", pageNo) // -> "event:1"
	cached, errcach := handler.Cache.Client.Get(ctx, key).Result()
	if errcach != redis.Nil {
		var events []*models.User
		fmt.Println("Cached wala data aaya :)")
		json.Unmarshal([]byte(cached), &events)
		c.JSON(http.StatusOK, gin.H{"data": events, "cached": true})
		return
	}
	fmt.Println("Cached data nhi aaya :(")
	allUsers, err := handler.UserHandler.GetAllUsersservice(pageNo)
	if len(allUsers) == 0 {
		response.ErrorResponse(c, http.StatusNotFound, "No Results found for this page")
	}
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	//Storing the cached data
	data, _ := json.Marshal(allUsers)
	handler.Cache.Client.Set(ctx, key, data, 2*time.Minute)
	c.JSON(http.StatusOK, gin.H{"cached": false, "data": allUsers})
}

//Get Full data on Attendees

func (handler *HandlerStruct) GetAllAttendeesList(c *gin.Context) {
	ctx := context.Background()
	key := "getAllattendees" // -> "event:1"
	cached, errcach := handler.Cache.Client.Get(ctx, key).Result()
	if errcach != redis.Nil {
		var event []*models.Attendees
		json.Unmarshal([]byte(cached), &event)
		c.JSON(http.StatusOK, gin.H{"data": event, "cached": true})
		return
	}

	allAttendees, err := handler.AttendeeHandler.GetAllAttendeesandEventList()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	data, _ := json.Marshal(allAttendees)
	handler.Cache.Client.Set(ctx, key, data, 2*time.Minute)
	response.SuccessResponse(c, allAttendees)
}

func (handler *HandlerStruct) RegisterUser(c *gin.Context) {
	var register *models.RegisterRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)

	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	register.Password = string(hashedPass)

	user := &models.User{
		Email:    register.Email,
		Password: register.Password,
		Name:     register.Name,
	}

	err3 := handler.UserHandler.Insertservice(user)

	if err3 != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err3.Error())
	}
	handler.Hub.BroadcastJSON(user)
	response.Response(c, http.StatusCreated, user)
}

func (handler *HandlerStruct) LoginUser(c *gin.Context) {
	var loginUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON")
		return
	}
	user, err := handler.UserHandler.FindByEmailPassword(loginUser.Email)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		response.ErrorResponse(c, http.StatusUnauthorized, "User Not Found")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password)) != nil {
		response.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	//Sessions Creation Start
	session := sessions.Default(c)
	session.Set("userID", user.Id)
	session.Set("userName", user.Name)
	session.Set("userEmail", user.Email)
	if err := session.Save(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Error Fetching Session")
		return
	}
	//Sessions Creation End
	rolesList, err := handler.UserRoleHandler.GetRolesByUser(user.Id)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Error Getting Roles from User")
		return
	}
	token, _ := utils.GenerateJWTRSA(user.Id, rolesList)
	response.SuccessResponse(c, gin.H{"Token": token, "Logged Status": "Logged In"})
}
