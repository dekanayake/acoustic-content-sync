package api

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"sync"
)

var once sync.Once

var instance *resty.Client

func Connect() *resty.Client {
	once.Do(func() {
		authUserName := env.AcousticAuthUserName()
		password := env.AcousticAuthPassword()
		apiKey := env.AcousticAPIKey()
		if authUserName == "" && apiKey == "" {
			log.Panic("No either user name of api values is provided ")
		}
		if authUserName != "" {
			if password == "" {
				log.Panic("Password not provided for acoustic user auth for user name :" + authUserName)
			}
			log.WithField("User name", authUserName).Info("Setting the user name as basic auth")
			instance = resty.New().SetBasicAuth(authUserName, password).SetDebug(env.IsDebugEnabled())
		} else if apiKey != "" {
			log.WithField("APIKey", apiKey).Info("Setting the api key as basic auth")
			instance = resty.New().SetBasicAuth("AcousticAPIKey", apiKey).SetDebug(env.IsDebugEnabled())
		}
	})
	return instance
}
