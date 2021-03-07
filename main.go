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

func nonreece(feedName string, configName string) {

	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	//catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	//err := catService.Create("Moodboard Non Reece Brands", feedName, configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = catService.Create("Moodboard Non Reece Categories", feedName, configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = catService.Create("TRS Product Color", feedName, configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "ff8c36e0-cc3d-48a0-8efe-9a4de800ce14")
	status, err := contentService.Create("4c8b4730-7503-485a-9c8e-23af27c61307", feedName, configName)
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
	//err := deleteService.Delete("ff8c36e0-cc3d-48a0-8efe-9a4de800ce14", "Delete Non Reece Moodboard content", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = deleteService.Delete("ff8c36e0-cc3d-48a0-8efe-9a4de800ce14", "Delete Non Reece Moodboard images", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
}

func nonreeceswatches(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	err := catService.Create("Moodboard Non Reece Brands", feedName, configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

	err = catService.Create("Moodboard Non Reece Categories", feedName, configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

	err = catService.Create("TRS Product Color", feedName, configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "ff8c36e0-cc3d-48a0-8efe-9a4de800ce14")
	status, err := contentService.Create("4c8b4730-7503-485a-9c8e-23af27c61307", feedName, configName)
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

func badges(feedName string, configName string) {

	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err := deleteService.Delete("59f1a68b-f518-47a5-83d9-bd16e26c7daa", "Delete Ambulant badge data", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "59f1a68b-f518-47a5-83d9-bd16e26c7daa")
	status, err := contentService.Create("da86ef75-537b-4c67-aebf-e476e1d2a099", feedName, configName)
	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
	log.Info(" success created record count  :" + strconv.Itoa(len(status.Success)))
	if status.FailuresExist() {
		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
		status.PrintFailed()
	}
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

}

func createBrands() {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	err = catService.Create("Reece AU website product brand", "brands.csv", "config_reece_brands.yaml")
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}
}

func createCats() {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	err = catService.Create("Reece AU website category", "category.csv", "config_reece_cats.yaml")
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}
}

func deleteCats() {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	err = catService.Delete("Reece AU website category")
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}
}

func main() {
	//feedName := "ambulantAS14288.1.csv"
	////feedName := "ambulantAS14288.1_strp_down.csv"
	//configName := "config_ambulant_product_badges.yaml"
	//badges(feedName,configName)
	//deleteCats()
	//createCats()
	createBrands()
}
