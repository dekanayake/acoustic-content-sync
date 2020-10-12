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

func reece() {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	err = catService.Create("TRS Brands", "reece_full_products_output.csv", "config_reece_products.yaml")
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

	err = catService.Create("TRS Product Category", "reece_full_products_output.csv", "config_reece_products.yaml")
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

	err = catService.Create("TRS Product Color", "reece_full_products_output.csv", "config_reece_products.yaml")
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}
	//
	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "ff8c36e0-cc3d-48a0-8efe-9a4de800ce14")
	status, err := contentService.Create("f5fe4c5c-67db-465a-aba6-75618cdcbf30", "reece_full_products_output.csv", "config_reece_products.yaml")
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
	//err := deleteService.Delete("ff8c36e0-cc3d-48a0-8efe-9a4de800ce14", "Delete Reece Moodboard content", "config_reece_products.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = deleteService.Delete("ff8c36e0-cc3d-48a0-8efe-9a4de800ce14", "Delete Reece Moodboard images", "config_reece_products.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
}

func nonreece() {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	//catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	//err := catService.Create("Moodboard Non Reece Brands", "external_products_aura.csv", "config_non_reece_products.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = catService.Create("Moodboard Non Reece Categories", "external_products_aura.csv", "config_non_reece_products.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = catService.Create("TRS Product Color", "external_products_aura.csv", "config_non_reece_products.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "ff8c36e0-cc3d-48a0-8efe-9a4de800ce14")
	status, err := contentService.Create("4c8b4730-7503-485a-9c8e-23af27c61307", "external_products_aura.csv", "config_non_reece_products_aura.yaml")
	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
	log.Info(" success created record count  :" + strconv.Itoa(len(status.Success)))
	if status.FailuresExist() {
		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
		status.PrintFailed()
	}
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}
	//
	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err := deleteService.Delete("0e958689-13ce-4eda-9c68-bb4dcc09dd73", "Delete Non Reece Moodboard content", "config_non_reece_products.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err := deleteService.Delete("0e958689-13ce-4eda-9c68-bb4dcc09dd73", "Delete Non Reece Moodboard images", "config_non_reece_products.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
}

func nonreeceswatches() {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	err := catService.Create("Moodboard Non Reece Brands", "non_reece_caeserstone_swatches.csv", "config_non_reece_products_caeserstone_swatches.yaml")
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

	err = catService.Create("Moodboard Non Reece Categories", "non_reece_caeserstone_swatches.csv", "config_non_reece_products_caeserstone_swatches.yaml")
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

	err = catService.Create("TRS Product Color", "non_reece_caeserstone_swatches.csv", "config_non_reece_products_caeserstone_swatches.yaml")
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "ff8c36e0-cc3d-48a0-8efe-9a4de800ce14")
	status, err := contentService.Create("4c8b4730-7503-485a-9c8e-23af27c61307", "non_reece_caeserstone_swatches.csv", "config_non_reece_products_caeserstone_swatches.yaml")
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
	//err = deleteService.Delete("0e958689-13ce-4eda-9c68-bb4dcc09dd73", "Delete Non Reece Moodboard content", "config_non_reece_product_swatches.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = deleteService.Delete("0e958689-13ce-4eda-9c68-bb4dcc09dd73", "Delete Non Reece Moodboard images", "config_non_reece_product_swatches.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
}

func main() {

	nonreeceswatches()

}
