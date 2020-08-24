package api

import (
	"errors"
	log "github.com/sirupsen/logrus"
	koazee "github.com/wesovilabs/koazee"
	"gopkg.in/resty.v1"
	"net/http"
)

type AuthClient interface {
	Authenticate(userName string, password string) (*http.Cookie,error)
}

type authClient struct {
	c resty.Client
	acousticAuthUrl string
}

func NewAuthClient(acousticAuthUrl string) AuthClient {
	return &authClient{
		c: *resty.New(),
		acousticAuthUrl: acousticAuthUrl,
	}
}

func (authClient *authClient) Authenticate (userName string, password string) (*http.Cookie,error) {
	req := authClient.c.NewRequest().
		SetBasicAuth(userName,password)
	if resp, err := req.Get(authClient.acousticAuthUrl) ; err != nil {
		return nil,err
	}  else if resp.IsSuccess(){
		authTokenCookie :=  koazee.StreamOf(resp.Cookies()).
			Filter(func(cookie *http.Cookie) bool{
					return cookie.Name == "x-ibm-dx-user-auth"
			}).
			First().Val().(*http.Cookie)
			log.WithFields(log.Fields{
				"userName": userName,
			}).Info("Successfully authenticated user ")
		if authTokenCookie == nil {
			return nil,errors.New("no authentication token received  : ")
		}
		return authTokenCookie,nil
	} else {
		return nil,errors.New("error in authenticating  : " + resp.Status())
	}
}

