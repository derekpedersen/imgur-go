package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/derekpedersen/imgur-go/model"
	"github.com/jeanphorn/log4go"
)

// Credentials for the Imgur API
var credentials model.Authorization

func NewAlbumService(clientID string) {
	credentials.ClientID = clientID
}

func QueryAlbum(albumHash string) (json string, err error) {
	log4go.Info("Querying Album: %s", albumHash)

	url := "https://api.imgur.com/3/album/" + albumHash

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log4go.Error("Error creating request:\n %v", err)
		return "", err
	}

	req.Header.Add("authorization", "Client-ID "+credentials.ClientID)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log4go.Error("Error making request:\n %v", err)
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log4go.Error("Error reading res.Body:\n %v", err)
		return "", err
	}

	json = string(body)

	log4go.Info("Album JSON: %s", json)

	return json, nil
}

// GetAlbum
func GetAlbum(albumHash string) (*model.Album, error) {
	res := model.AlbumResponse{}
	albumJson, err := QueryAlbum(albumHash)
	if err != nil {
		log4go.Error("Error querying album:\n %v", err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(albumJson), &res); err != nil {
		log4go.Error("Error unmarshalling album response:\n %v", err)
		return nil, err
	}

	return &res.Data, nil
}
