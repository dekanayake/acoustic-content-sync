package main

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/csv"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {

	authClient := api.NewAuthClient("https://content-eu-1.content-cms.com/api/f43b0181-6ee9-4539-b3e9-f7f37a95f962/login/v1/basicauth")
	token,authErr  := authClient.Authenticate(os.Getenv("UserName"),os.Getenv("Password"))

	log.Error(authErr)

	service := csv.NewService(token,os.Getenv("ContentAuthUrl"),"0e958689-13ce-4eda-9c68-bb4dcc09dd73")
	err := service.Create("f5fe4c5c-67db-465a-aba6-75618cdcbf30","data.csv","config.yaml")

	log.Error(err)
}