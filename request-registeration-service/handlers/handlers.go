package handlers

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"request-registeration-service/configs"
	"request-registeration-service/messages"
	"request-registeration-service/models"
)

func SaveRequestHandler(ctx echo.Context) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	requestInfo := models.RequestInfo{
		Email:  req.Email,
		Status: "pending",
	}
	if err := configs.DB.Create(&requestInfo).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToSaveRequest)
	}
	err := ctx.Request().ParseMultipartForm(10 << 20)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	file, handler, err := ctx.Request().FormFile("file")
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	defer file.Close()
	fileBytes, err := os.ReadFile(handler.Filename)
	err = configs.UploadFile(bytes.NewReader(fileBytes), fmt.Sprintf("file-%d", requestInfo.ID))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToUploadFile)
	}
	return ctx.JSON(http.StatusOK, messages.RequestRegistered)
}
