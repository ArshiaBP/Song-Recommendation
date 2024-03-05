package handlers

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
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
		requestInfo.Status = "failure"
		configs.DB.Save(&requestInfo)
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	file, handler, err := ctx.Request().FormFile("file")
	if err != nil {
		requestInfo.Status = "failure"
		configs.DB.Save(&requestInfo)
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	defer file.Close()
	fileBytes, err := os.ReadFile(handler.Filename)
	err = configs.UploadFile(bytes.NewReader(fileBytes), fmt.Sprintf("file-%d", requestInfo.ID))
	if err != nil {
		requestInfo.Status = "failure"
		configs.DB.Save(&requestInfo)
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToUploadFile)
	}
	err = configs.Ch.PublishWithContext(configs.Ctx, "", configs.Queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(fmt.Sprint(requestInfo.ID)),
	})
	if err != nil {
		requestInfo.Status = "failure"
		configs.DB.Save(&requestInfo)
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToWriteInMQ)
	}
	return ctx.JSON(http.StatusOK, messages.RequestRegistered)
}
