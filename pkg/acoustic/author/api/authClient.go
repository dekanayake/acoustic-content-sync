package api

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"net/http"
)

type AuthClient interface {
	Authenticate(userName string, password string) ([]*http.Cookie, error)
}

type authClient struct {
	c               resty.Client
	acousticAuthUrl string
}

func NewAuthClient(acousticAuthUrl string) AuthClient {
	return &authClient{
		c:               *resty.New().SetDebug(env.IsDebugEnabled()),
		acousticAuthUrl: acousticAuthUrl,
	}
}

func (authClient *authClient) Authenticate(userName string, password string) ([]*http.Cookie, error) {
	req := authClient.c.NewRequest().
		SetBasicAuth(userName, password)
	if resp, err := req.Get(authClient.acousticAuthUrl); err != nil {
		return nil, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		log.WithFields(log.Fields{
			"userName": userName,
		}).Info("Successfully authenticated user ")
		return resp.Cookies(), nil
	} else {
		return nil, errors.ErrorMessageWithStack("error in authenticating  : " + resp.Status())
	}
}
