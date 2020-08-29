package main

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/csv"
	"github.com/dekanayake/acoustic-content-sync/pkg/logrus"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	errorHandling := logrus.PkgErrorEntry{Entry: log.New().WithField("", "")}

	service := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	err := service.Create("TRS Brands", "output.csv", "config.yaml")
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

	//service := csv.NewContentService(os.Getenv("AcousticAPIURL"), "0e958689-13ce-4eda-9c68-bb4dcc09dd73")
	//err := service.Create("f5fe4c5c-67db-465a-aba6-75618cdcbf30", "data.csv", "config.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
}
