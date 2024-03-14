package services

import (
	"log"
	"songID-identification-service/configs"
)

func saveSpotifyID(requestID int) error {
	songID, err := getSongID(requestID)
	if err != nil {
		return err
	}
	log.Println("spotify api worked successfully")
	result := configs.DB.Table("request_infos").Where("id = ?", requestID).Update("song_id", songID)
	if result.Error != nil {
		return err
	}
	result = configs.DB.Table("request_infos").Where("id = ?", requestID).Update("status", "ready")
	return result.Error
}
