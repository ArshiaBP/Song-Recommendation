package services

import "email-service/configs"

func checkStatus() (string, string, string, error) {
	var infos struct {
		email  string
		songID string
		status string
	}
	result := configs.DB.Table("request_infos").Where("status in ?", []string{"ready", "failure"}).Select([]string{"email", "song_id", "status"}).Scan(&infos)
	if result.Error != nil {
		return "", "", "", result.Error
	}
	return infos.email, infos.songID, infos.status, nil
}
