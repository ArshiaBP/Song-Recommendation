package services

import (
	"email-service/configs"
	"email-service/models"
)

func CheckStatus() ([]models.RequestInfo, error) {
	var requests []models.RequestInfo
	result := configs.DB.Select([]string{"email", "song_id", "status"}).Where("status in ?", []string{"ready", "failure"}).Find(&requests)
	if result.Error != nil {
		return []models.RequestInfo{}, result.Error
	}
	return requests, nil
}
