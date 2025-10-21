package handlers

import (
	"Gin_Event_App_Arc/fabric"
	"Gin_Event_App_Arc/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (handler *HandlerStruct) GetUserFabric(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid User ID in Request Params!")
		return
	}
	userFabric, err2 := handler.UserFabricHandler.GetUserFabricService(id)
	if err2 != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err2.Error())
		return
	}
	response.Response(c, http.StatusOK, userFabric)
}

func (handler *HandlerStruct) UpdateUserFabric(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid User ID in Request Params!")
		return
	}
	userFabric := &fabric.UserUpdate{}                    //Ek nya Struct ka instance create kiya for User in Fabric
	if err := c.ShouldBindJSON(&userFabric); err != nil { //Bind from the body
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	errNow := handler.UserFabricHandler.UpdateUserFabricService(*userFabric, id)
	if errNow != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, errNow.Error())
		return
	}
	response.Response(c, http.StatusOK, nil)
}

func (handler *HandlerStruct) DeleteUserFabric(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid User ID in Request Params!")
		return
	}
	errNow := handler.UserFabricHandler.DeleteUserFabricService(id)
	if errNow != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Error deleting the User from fabric")
		return
	}
	response.Response(c, http.StatusNoContent, errNow)
}

func (handler *HandlerStruct) InsertUserFabric(c *gin.Context) {
	userFabric := &fabric.User{}                          //Ek nya Struct ka instance create kiya for User in Fabric
	if err := c.ShouldBindJSON(&userFabric); err != nil { //Bind from the body
		fmt.Println("1")
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	errNow := handler.UserFabricHandler.InsertUserFabricService(*userFabric)
	if errNow != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, errNow.Error())
		return
	}
	response.Response(c, http.StatusCreated, nil)
}

func (handler *HandlerStruct) GetAllUsersFabric(c *gin.Context) {
	result, err := handler.UserFabricHandler.GetAllUsersFabricService()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("WTHHHH")
	response.Response(c, http.StatusOK, result)
}

func (handler *HandlerStruct) GetFullHistory(c *gin.Context) {
	result, err := handler.UserFabricHandler.GetFullHistoryService()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Response(c, http.StatusOK, result)
}

func (handler *HandlerStruct) InitLedger(c *gin.Context) {
	result := handler.UserFabricHandler.InitLedgerService()
	response.Response(c, http.StatusOK, result)
}
