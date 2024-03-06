package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"songID-identification-service/configs"
	"strconv"
	"strings"
)

const shazamURL = "https://shazam-api-free.p.rapidapi.com/shazam/recognize/"
const spotifyBaseURL = "https://spotify23.p.rapidapi.com/search/?q="

func Consume() {
	messages, err := configs.Ch.Consume(configs.Queue.Name, "", false, false, false, false, nil)
	if err != nil {
		Consume()
	}
	quit := make(chan struct{})
	go func() {
		for msg := range messages {
			requestID, _ := strconv.Atoi(string(msg.Body))
			msg.Ack(true)
			go func() {
				err = saveSpotifyID(requestID)
				if err != nil {
					configs.DB.Table("request_infos").Where("id = ?", requestID).Update("status", "failure")
				}
			}()
		}
	}()
	<-quit
}

func getSong(requestID int) ([]byte, error) {
	fileName := fmt.Sprintf("file-%d", requestID)
	fileBytes, err := configs.DownloadFile(fileName)
	if err != nil {
		return []byte{}, err
	}
	return fileBytes, nil
}

func getSongTitle(requestID int) (string, error) {
	var track struct {
		title string
	}
	fileBytes, err := getSong(requestID)
	if err != nil {
		return "", err
	}
	payload := bytes.NewReader(fileBytes)
	req, _ := http.NewRequest("POST", shazamURL, payload)
	req.Header.Add("content-type", "multipart/form-data; boundary=---011000010111000001101001")
	req.Header.Add("X-RapidAPI-Key", os.Getenv("API_Key"))
	req.Header.Add("X-RapidAPI-Host", "shazam-api-free.p.rapidapi.com")
	response, _ := http.DefaultClient.Do(req)
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&track)
	if err != nil {
		return "", err
	}
	return track.title, nil
}

func getSongID(requestID int) (string, error) {
	var tracks struct {
		items struct {
			data struct {
				id string
			}
		}
	}
	title, err := getSongTitle(requestID)
	if err != nil || title == "" {
		return "", err
	}
	titleWords := strings.Split(title, " ")
	spotifyURL := fmt.Sprintf("%s%s", spotifyBaseURL, titleWords[0])
	for i, word := range titleWords {
		if i == 0 {
			continue
		}
		spotifyURL = fmt.Sprintf(spotifyURL + "%20" + word)
	}
	spotifyURL = fmt.Sprintf("%s%s", spotifyURL, "&type=tracks&offset=0&limit=10&numberOfTopResults=1")
	req, _ := http.NewRequest("GET", spotifyURL, nil)
	req.Header.Add("X-RapidAPI-Key", os.Getenv("API_Key"))
	req.Header.Add("X-RapidAPI-Host", "spotify23.p.rapidapi.com")
	response, _ := http.DefaultClient.Do(req)
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&tracks)
	if err != nil {
		return "", err
	}
	return tracks.items.data.id, nil
}

func saveSpotifyID(requestID int) error {
	songID, err := getSongID(requestID)
	if err != nil {
		return err
	}
	result := configs.DB.Table("request_infos").Where("id = ?", requestID).Update("song_id", songID)
	if result.Error != nil {
		return err
	}
	result = configs.DB.Table("request_infos").Where("id = ?", requestID).Update("status", "ready")
	return result.Error
}
