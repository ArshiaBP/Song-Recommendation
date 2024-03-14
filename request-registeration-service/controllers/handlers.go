package controllers

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"io"
	"net/http"
	"request-registeration-service/configs"
	"request-registeration-service/messages"
	"request-registeration-service/models"
)

func SaveRequestHandler(ctx echo.Context) error {
	email := ctx.Param("email")
	requestInfo := models.RequestInfo{
		Email:  email,
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
	file, _, err := ctx.Request().FormFile("file")
	if err != nil {
		requestInfo.Status = "failure"
		configs.DB.Save(&requestInfo)
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		requestInfo.Status = "failure"
		configs.DB.Save(&requestInfo)
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	err = configs.UploadFile(bytes.NewReader(fileBytes), fmt.Sprintf("file-%d.mp3", requestInfo.ID))
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
