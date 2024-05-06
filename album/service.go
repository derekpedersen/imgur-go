package album

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/derekpedersen/imgur-go/authorization"
	log "github.com/sirupsen/logrus"
)

// AlbumService interface
type AlbumService interface {
	QueryAlbum(albumHash string) (json *string, err error)
	GetAlbum(albumHash string) (*Album, error)
}

// AlbumServiceImpl struct
type AlbumServiceImpl struct {
	auth authorization.Authorization
	url  string
}

// NewAlbumService creates a new album service
func NewAlbumService(
	auth authorization.Authorization,
	apiURL string,
) *AlbumServiceImpl {
	return &AlbumServiceImpl{
		auth: auth,
		url:  apiURL,
	}
}

// QueryAlbum queries an album
func (svc *AlbumServiceImpl) QueryAlbum(
	albumHash string,
) (
	json *string,
	err error,
) {

	log.Infof("Querying Album: %s", albumHash)

	req, err := http.NewRequest("GET", svc.url+albumHash, nil)
	if err != nil {
		log.Errorf("Error creating request:\n %v", err)
		return nil, err
	}

	req.Header.Add("Authorization", " Bearer "+svc.auth.ImgurTokenResponse.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Error making request:\n %v", err)
		return nil, err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error reading res.Body:\n %v", err)
		return nil, err
	}

	body := string(b)

	json = &body

	log.Debugf("Album JSON: %v", *json)

	return json, nil
}

// GetAlbum gets the album
func (svc *AlbumServiceImpl) GetAlbum(albumHash string) (*Album, error) {
	res := AlbumResponse{}
	albumJSON, err := svc.QueryAlbum(albumHash)
	if err != nil {
		log.Errorf("Error querying album:\n %v", err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(*albumJSON), &res); err != nil {
		log.Errorf("Error unmarshalling album response:\n %v", err)
		return nil, err
	}

	return &res.Data, nil
}
