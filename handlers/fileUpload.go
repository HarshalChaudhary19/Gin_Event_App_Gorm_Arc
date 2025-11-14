package handlers

import (
	"Gin_Event_App_Arc/models"
	"Gin_Event_App_Arc/response"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (handler *HandlerStruct) UploadFile(c *gin.Context) {

	//Getting the file
	fh, err := c.FormFile("file")
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Cannot Read File")
		return
	}
	//Opening the file
	f, err := fh.Open()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer f.Close()
	//Reading file
	data, err := io.ReadAll(f)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	contentType := fh.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	file := models.FileUpload{
		FileName:    fh.Filename,
		ContentType: contentType,
		Data:        data,
		Size:        len(data),
	}
	errnew := handler.FileUploadHandler.UploadFile(&file)
	if errnew != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Response(c, http.StatusCreated, file)
}

func (handler *HandlerStruct) UpdateFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Cannot Read File")
		return
	}
	//Opening the file
	f, err := fh.Open()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer f.Close()
	//Reading file
	data, err := io.ReadAll(f)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	contentType := fh.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	file := models.FileUpload{
		ID:          uint(id),
		FileName:    fh.Filename,
		ContentType: contentType,
		Data:        data,
		Size:        len(data),
	}
	fileReturned, err := handler.FileUploadHandler.UpdateFile(&file)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	disposition := "inline"
	if c.Query("download") == "true" {
		disposition = "attachment"
	}
	c.Header("Content-Type", fileReturned.ContentType)
	c.Header("Content-Disposition", disposition+"; filename=\""+fileReturned.FileName+"\"")
	c.Data(http.StatusOK, fileReturned.ContentType, fileReturned.Data)
}

func (handler *HandlerStruct) DeleteFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = handler.FileUploadHandler.DeleteFile(id)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Response(c, http.StatusNoContent, "Done")
}

func (handler *HandlerStruct) GetFileByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	file, err := handler.FileUploadHandler.GetFile(id)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	disposition := "inline"
	if c.Query("download") == "true" {
		disposition = "attachment"
	}
	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Disposition", disposition+"; filename=\""+file.FileName+"\"")
	c.Data(http.StatusOK, file.ContentType, file.Data)
}

func (handler *HandlerStruct) GetAllFiles(c *gin.Context) {
	// files, err := handler.FileUploadHandler.GetAllFiles()
	// if err != nil {
	// 	response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch files")
	// 	return
	// }

}
