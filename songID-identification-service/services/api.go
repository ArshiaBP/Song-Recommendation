package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	shazamURL      = "https://shazam-api-free.p.rapidapi.com/shazam/recognize/"
	spotifyBaseURL = "https://spotify23.p.rapidapi.com/search/?q="
	shazamHost     = "shazam-api-free.p.rapidapi.com"
	spotifyHost    = "spotify23.p.rapidapi.com"
)

func getSongTitle(requestID int) (string, error) {
	type Track struct {
		Title string `json:"title"`
	}
	var res struct {
		Track Track `json:"track"`
	}
	var (
		buf = new(bytes.Buffer)
		w   = multipart.NewWriter(buf)
	)
	fileName := fmt.Sprintf("file-%d.mp3", requestID)
	fieldName := "upload_file"
	fileBytes, err := getSong(requestID)
	if err != nil {
		log.Println("download file failed")
		return "", err
	}
	log.Println("file downloaded")
	part, err := w.CreateFormFile(fieldName, filepath.Base(fileName))
	if err != nil {
		return "", err
	}
	_, err = part.Write(fileBytes)
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", shazamURL, buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+w.Boundary())
	req.Header.Add("X-RapidAPI-Key", os.Getenv("API_Key"))
	req.Header.Add("X-RapidAPI-Host", shazamHost)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	resBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return "", err
	}
	return res.Track.Title, nil
}

func getSongID(requestID int) (string, error) {
	type Data struct {
		ID string `json:"id"`
	}
	type Items struct {
		Data Data `json:"data"`
	}
	type Tracks struct {
		Items []Items `json:"items"`
	}
	var res struct {
		Tracks Tracks `json:"tracks"`
	}
	title, err := getSongTitle(requestID)
	if err != nil {
		log.Println("shazam api failed")
		return "", err
	}
	if title == "" {
		log.Println("shazam api failed")
		return "", errors.New("empty title")
	}
	log.Println("shazam api worked successfully")
	titleWords := strings.Split(title, " ")
	spotifyURL := fmt.Sprintf("%s%s", spotifyBaseURL, titleWords[0])
	for i, word := range titleWords {
		if i == 0 {
			continue
		}
		temp := "%20"
		spotifyURL = fmt.Sprintf("%s%s%s", spotifyURL, temp, word)
	}
	spotifyURL = fmt.Sprintf("%s%s", spotifyURL, "&type=tracks&offset=0&limit=1&numberOfTopResults=1")
	req, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-RapidAPI-Key", os.Getenv("API_Key"))
	req.Header.Add("X-RapidAPI-Host", spotifyHost)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("spotify api failed")
		return "", err
	}
	resBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("spotify api failed")
		return "", err
	}
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		log.Println("spotify api failed")
		return "", err
	}
	return res.Tracks.Items[0].Data.ID, nil
}
