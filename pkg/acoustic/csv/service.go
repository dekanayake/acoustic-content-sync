package csv

import (
	"github.com/wesovilabs/koazee"
	"net/http"
)

type Service interface {
	Create(contentType string,dataFeedPath string, configPath string) error
}

type service struct {
	authToken *http.Cookie
	acousticAuthApiUrl string
	acousticContentLib string
}




func NewService(authToken *http.Cookie, acousticAuthApiUrl string, acousticContentLib string) Service {
	return &service{
		authToken : authToken,
		acousticAuthApiUrl :acousticAuthApiUrl,
		acousticContentLib:acousticContentLib,
	}
}


func (service *service) Create(contentType string, dataFeedPath string, configPath string) error {



}


