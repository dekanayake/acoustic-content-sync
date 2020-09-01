package main

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/csv"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/logrus"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

func init() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	acousticSyncLog, err := os.OpenFile("acoustic_sync.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Panic("Error in writing to log file", err)
	}
	log.SetFormatter(Formatter)
	mw := io.MultiWriter(os.Stdout, acousticSyncLog)
	log.SetOutput(mw)
	log.SetLevel(log.InfoLevel)
}

func main() {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}

	//catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	//err := catService.Create("TRS Brands", "data.csv", "config.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = catService.Create("TRS Product Category", "data.csv", "config.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = catService.Create("TRS Product Color", "data.csv", "config.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "0e958689-13ce-4eda-9c68-bb4dcc09dd73")
	status, err := contentService.Create("f5fe4c5c-67db-465a-aba6-75618cdcbf30", "data.csv", "config.yaml")
	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
	log.Info(" success created record count  :" + strconv.Itoa(len(status.Success)))
	if status.FailuresExist() {
		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
		status.PrintFailed()
	}
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err := deleteService.Delete("0e958689-13ce-4eda-9c68-bb4dcc09dd73", "Delete Reece Moodboard content", "config.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = deleteService.Delete("0e958689-13ce-4eda-9c68-bb4dcc09dd73", "Delete Reece Moodboard images", "config.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

}
