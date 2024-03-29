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

func reece(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	//catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	//err = catService.CreateBatch("TRS Brands", "reece_full_products_output.csv", "config_reece_products.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = catService.CreateBatch("TRS Product Category", "reece_full_products_output.csv", "config_reece_products.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = catService.CreateBatch("TRS Product Color", "reece_full_products_output.csv", "config_reece_products.yaml")
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "ff8c36e0-cc3d-48a0-8efe-9a4de800ce14")
	status, err := contentService.CreateBatch("f5fe4c5c-67db-465a-aba6-75618cdcbf30", feedName, configName)
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

func delete_reece(configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	err = deleteService.Delete("ff8c36e0-cc3d-48a0-8efe-9a4de800ce14", "Delete Reece Moodboard content", configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

	err = deleteService.Delete("ff8c36e0-cc3d-48a0-8efe-9a4de800ce14", "Delete Reece Moodboard images", configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}
}

func nonreece(feedName string, configName string) {

	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	//catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	//err := catService.CreateBatch("Moodboard Non Reece Brands", feedName, configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = catService.CreateBatch("Moodboard Non Reece Categories", feedName, configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	//err = catService.CreateBatch("TRS Product Color", feedName, configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "ff8c36e0-cc3d-48a0-8efe-9a4de800ce14")
	status, err := contentService.CreateBatch("4c8b4730-7503-485a-9c8e-23af27c61307", feedName, configName)
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
	status, err := contentService.CreateBatch("4c8b4730-7503-485a-9c8e-23af27c61307", feedName, configName)
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
	status, err := contentService.CreateBatch("da86ef75-537b-4c67-aebf-e476e1d2a099", feedName, configName)
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

func sliBanners(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err := deleteService.Delete("24f48506-0498-41f3-8701-43ab4a81b396", "Delete Sli banners", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//err = deleteService.Delete("24f48506-0498-41f3-8701-43ab4a81b396", "Delete Sli banner images", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "24f48506-0498-41f3-8701-43ab4a81b396")
	status, err := contentService.CreateBatch("92423321-23ce-423e-bff2-74d39a99e449", feedName, configName)
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

func brands(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err := deleteService.Delete("191c5ef1-9194-4326-bc4f-bc37d681685e", "Delete brands", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "191c5ef1-9194-4326-bc4f-bc37d681685e")
	status, err := contentService.CreateBatch("9c4cef56-78fc-41d1-b959-d5ddb8f4fd9d", feedName, configName)
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

func customSearch(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err := deleteService.Delete("191c5ef1-9194-4326-bc4f-bc37d681685e", "Delete brands", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "970266a6-c90c-4cc9-8bd6-e5867920e5ce")
	status, err := contentService.CreateBatch("9c95a805-99a2-41a4-b423-0bdfa592cc53", feedName, configName)
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

func createCats(catName string, feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	err = catService.Create(catName, feedName, configName)
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
	feedName := "customorder_nz.csv"
	configName := "config_customorder_nz_product_badges.yaml"
	//delete_reece(configName)
	badges(feedName, configName)
}
