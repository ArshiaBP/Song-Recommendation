package services

import (
	"fmt"
	"songID-identification-service/configs"
)

func getSong(requestID int) ([]byte, error) {
	fileName := fmt.Sprintf("file-%d.mp3", requestID)
	fileBytes, err := configs.DownloadFile(fileName)
	if err != nil {
		return []byte{}, err
	}
	return fileBytes, nil
}
