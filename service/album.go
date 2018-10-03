package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/derekpedersen/imgur-go/model"
	log "github.com/sirupsen/logrus"
)

// AlbumService interface
type AlbumService interface {
	QueryAlbum(albumHash string) (json string, err error)
	GetAlbum(albumHash string) (*model.Album, error)
}

// AlbumServiceImpl struct
type AlbumServiceImpl struct {
	auth model.Authorization
}

// NewAlbumService creates a new album service
func NewAlbumService(clientID string) *AlbumServiceImpl {
	return &AlbumServiceImpl{
		auth: model.Authorization{
			ClientID: clientID,
		},
	}
}

// QueryAlbum queries an album
func (svc *AlbumServiceImpl) QueryAlbum(albumHash string) (json string, err error) {
	log.Infof("Querying Album: %s", albumHash)

	url := "https://api.imgur.com/3/album/" + albumHash

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Error creating request:\n %v", err)
		return "", err
	}

	req.Header.Add("authorization", "Client-ID "+svc.auth.ClientID)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Error making request:\n %v", err)
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error reading res.Body:\n %v", err)
		return "", err
	}

	json = string(body)

	log.Debugf("Album JSON: %s", json)

	return json, nil
}

// GetAlbum gets the album
func (svc *AlbumServiceImpl) GetAlbum(albumHash string) (*model.Album, error) {
	res := model.AlbumResponse{}
	albumJSON, err := svc.QueryAlbum(albumHash)
	if err != nil {
		log.Errorf("Error querying album:\n %v", err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(albumJSON), &res); err != nil {
		log.Errorf("Error unmarshalling album response:\n %v", err)
		return nil, err
	}

	return &res.Data, nil
}
