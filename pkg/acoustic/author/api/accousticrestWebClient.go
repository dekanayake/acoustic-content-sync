package api

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"gopkg.in/resty.v1"
	"log"
	"sync"
)

var once sync.Once

var instance *resty.Client

func Connect() *resty.Client {
	once.Do(func() {
		token, err := NewAuthClient(env.AcousticAuthUrl()).Authenticate(env.AcousticAuthUserName(), env.AcousticAuthPassword())
		if err != nil {
			log.Panic("auth failed", err)
		}
		instance = resty.New().SetCookie(token).SetDebug(env.IsDebugEnabled())
	})

	return instance
}
