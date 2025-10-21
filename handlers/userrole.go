package handlers

import (
	"Gin_Event_App_Arc/response"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (handler *HandlerStruct) UpdateUserRole(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid User ID in Request Params!")
		return
	}

	role_id, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid Role ID in Request Params!")
		return
	}
	err = handler.UserRoleHandler.Update(user_id, role_id)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}

func (handler *HandlerStruct) GetRolesByUser(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid User ID in Request Params!")
		return
	}
	rolesList, err := handler.UserRoleHandler.GetRolesByUser(user_id)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if rolesList == nil {
		response.ErrorResponse(c, http.StatusNotFound, "For This userId there is no Role")
		return
	}
	response.SuccessResponse(c, rolesList)
}

func (handler *HandlerStruct) GetUsersByRole(c *gin.Context) {
	role_id, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid User ID in Request Params!")
		return
	}
	userList, err := handler.UserRoleHandler.GetUsersByRole(role_id)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if userList == nil {
		response.ErrorResponse(c, http.StatusNotFound, "No User for this Role")
		return
	}
	response.SuccessResponse(c, userList)
}

func (handler *HandlerStruct) GetAll(c *gin.Context) {
	alluserRoleList, err := handler.UserRoleHandler.GetAll()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessResponse(c, alluserRoleList)
}

func logoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	// remove session data and expire cookie
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1}) // instruct cookie deletion
	_ = session.Save()
	response.SuccessResponse(c, "logged out")
}
