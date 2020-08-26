package main

import (
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




	service := csv.NewService(os.Getenv("AcousticAPIURL"),"0e958689-13ce-4eda-9c68-bb4dcc09dd73")
	err := service.Create("f5fe4c5c-67db-465a-aba6-75618cdcbf30","data.csv","config.yaml")

	log.Error(err)
}