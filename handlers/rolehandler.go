package handlers

import (
	"Gin_Event_App_Arc/models"
	"Gin_Event_App_Arc/response"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func (handler *HandlerStruct) CreateRole(c *gin.Context) {
	var role models.Role

	if err := c.ShouldBindJSON(&role); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := handler.RoleHandler.Insert(&role) //Calling to the service

	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Response(c, http.StatusCreated, role)
}

func (handler *HandlerStruct) UpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid Role ID in request params!")
		return
	}
	existingRole, err := handler.RoleHandler.Get(id)

	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	if existingRole == nil {
		response.ErrorResponse(c, http.StatusNotFound, "No such Event Exists ")
		return
	}

	updatedRole := &models.Role{} //Ek nya Struct ka instance create kiya

	if err := c.ShouldBindJSON(&updatedRole); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	updatedRole.Id = id
	updatedRolen, err2 := handler.RoleHandler.Update(updatedRole, id)

	if err2 != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err2.Error())
		return
	}
	response.SuccessResponse(c, updatedRolen)
}

func (handler *HandlerStruct) DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid Role ID in request params!")
		return
	}
	if err := handler.EventHandler.Delete(id); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	response.Response(c, http.StatusNoContent, nil)
}

func (handler *HandlerStruct) GetRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid Role ID in request params!")
		return
	}

	//Caching check
	ctx := context.Background()
	key := fmt.Sprintf("RoleByID:%d", id) // -> "event:1"
	cached, errcach := handler.Cache.Client.Get(ctx, key).Result()
	if errcach != redis.Nil {
		var role *models.Role
		json.Unmarshal([]byte(cached), &role)
		c.JSON(http.StatusOK, gin.H{"data": role, "cached": true})
		return
	}

	roleFromDB, err1 := handler.RoleHandler.Get(id)

	if roleFromDB == nil {
		response.ErrorResponse(c, http.StatusNotFound, "No such Role Exists")
		return
	}

	if reflect.DeepEqual(roleFromDB, models.Role{}) {
		response.ErrorResponse(c, http.StatusNotFound, "No such Event Exists ")
	}

	if err1 != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err1.Error())
	}

	data, _ := json.Marshal(roleFromDB)
	handler.Cache.Client.Set(ctx, key, data, 2*time.Minute)
	response.SuccessResponse(c, roleFromDB)
}

func (handler *HandlerStruct) GetAllRole(c *gin.Context) {
	allRoles, err := handler.RoleHandler.GetAll()
	if len(allRoles) == 0 {
		response.ErrorResponse(c, http.StatusNotFound, "Page Not Found")
		return
	}
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Error fetching List of Roles")
		return
	}
	response.SuccessResponse(c, allRoles)
}
