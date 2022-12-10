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

//
//func reece(feedName string, configName string) {
//	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
//	var err error
//	//catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
//	//err = catService.CreateBatch("TRS Brands", "reece_full_products_output.csv", "config_reece_products.yaml")
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//	//
//	//err = catService.CreateBatch("TRS Product Category", "reece_full_products_output.csv", "config_reece_products.yaml")
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//	//
//	//err = catService.CreateBatch("TRS Product Color", "reece_full_products_output.csv", "config_reece_products.yaml")
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//	//
//	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "ff8c36e0-cc3d-48a0-8efe-9a4de800ce14")
//	status, err := contentService.CreateBatch("f5fe4c5c-67db-465a-aba6-75618cdcbf30", feedName, configName)
//	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
//	log.Info(" success created record count  :" + strconv.Itoa(len(status.Success)))
//	if status.FailuresExist() {
//		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
//		status.PrintFailed()
//	}
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//
//	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
//	//err := deleteService.Delete("ff8c36e0-cc3d-48a0-8efe-9a4de800ce14", "Delete Reece Moodboard content", "config_reece_products.yaml")
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//	//
//	//err = deleteService.Delete("ff8c36e0-cc3d-48a0-8efe-9a4de800ce14", "Delete Reece Moodboard images", "config_reece_products.yaml")
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//}

func threed_materials(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner materials", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner materials glb", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "cbcf8e69-a6ad-474c-b5db-766fc6956cfb")
	status, err := contentService.CreateBatch("07198b81-3dc0-4131-9f4c-17f3d1640049", feedName, configName)
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

func threed_materials_update(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner materials", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner materials glb", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "cbcf8e69-a6ad-474c-b5db-766fc6956cfb")
	status, err := contentService.CreateBatch("07198b81-3dc0-4131-9f4c-17f3d1640049", feedName, configName)
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

func threed_materials_configuration(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	//
	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner configuration", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "cbcf8e69-a6ad-474c-b5db-766fc6956cfb")
	status, err := contentService.CreateBatch("6972fdb3-f431-4a8b-81f5-95b34148f616", feedName, configName)
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

func threed_paint(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	//
	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner configuration", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "cbcf8e69-a6ad-474c-b5db-766fc6956cfb")
	status, err := contentService.CreateBatch("6a53e219-a702-4f4d-a8f2-31d1f9176cb8", feedName, configName)
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

func threed_tiles_update(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	//
	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner configuration", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "cbcf8e69-a6ad-474c-b5db-766fc6956cfb")
	status, err := contentService.CreateBatch("ef48246d-a3e4-4cf7-aee7-32a19d9afd20", feedName, configName)
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

func threed_materials_configuration_update(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	//
	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner configuration", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "cbcf8e69-a6ad-474c-b5db-766fc6956cfb")
	status, err := contentService.CreateBatch("6972fdb3-f431-4a8b-81f5-95b34148f616", feedName, configName)
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

func threed_materials_review(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner materials", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner materials glb", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "af32c812-305d-4cea-9136-dd5becc10b1d")
	status, err := contentService.CreateBatch("07198b81-3dc0-4131-9f4c-17f3d1640049", feedName, configName)
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

func threed_materials_configuration_review(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	//
	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner configuration", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "af32c812-305d-4cea-9136-dd5becc10b1d")
	status, err := contentService.CreateBatch("6972fdb3-f431-4a8b-81f5-95b34148f616", feedName, configName)
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

func custom_search(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	//
	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner configuration", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "970266a6-c90c-4cc9-8bd6-e5867920e5ce")
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

func custom_search_products(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	//
	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err = deleteService.Delete("cbcf8e69-a6ad-474c-b5db-766fc6956cfb", "Delete 3dplanner configuration", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "970266a6-c90c-4cc9-8bd6-e5867920e5ce")
	status, err := contentService.CreateBatch("945ba2f8-5393-41f3-8504-7ea6e5086ed8", feedName, configName)
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

func delete_pig_all(configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err = deleteService.Delete("f12876da-0017-41a8-8008-65065ed59f8a", "Delete Project Inspiration Gallery content", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
	//
	err = deleteService.Delete("f12876da-0017-41a8-8008-65065ed59f8a", "Delete Project Inspiration Teams", configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

	//err = deleteService.Delete("f12876da-0017-41a8-8008-65065ed59f8a", "Delete Project Inspiration Gallery  images", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}
}

//
//func delete_reece(configName string) {
//	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
//	var err error
//	deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
//	err = deleteService.Delete("ff8c36e0-cc3d-48a0-8efe-9a4de800ce14", "Delete Reece Moodboard content", configName)
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//
//	err = deleteService.Delete("ff8c36e0-cc3d-48a0-8efe-9a4de800ce14", "Delete Reece Moodboard images", configName)
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//}
//
//func nonreece(feedName string, configName string) {
//
//	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
//	//catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
//	//err := catService.CreateBatch("Moodboard Non Reece Brands", feedName, configName)
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//	//
//	//err = catService.CreateBatch("Moodboard Non Reece Categories", feedName, configName)
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//	//
//	//err = catService.CreateBatch("TRS Product Color", feedName, configName)
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//
//	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "ff8c36e0-cc3d-48a0-8efe-9a4de800ce14")
//	status, err := contentService.CreateBatch("4c8b4730-7503-485a-9c8e-23af27c61307", feedName, configName)
//	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
//	log.Info(" success created record count  :" + strconv.Itoa(len(status.Success)))
//	if status.FailuresExist() {
//		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
//		status.PrintFailed()
//	}
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//
//	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
//	//err := deleteService.Delete("ff8c36e0-cc3d-48a0-8efe-9a4de800ce14", "Delete Non Reece Moodboard content", configName)
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//	//
//	//err = deleteService.Delete("ff8c36e0-cc3d-48a0-8efe-9a4de800ce14", "Delete Non Reece Moodboard images", configName)
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//}
//
//func nonreeceswatches(feedName string, configName string) {
//	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
//	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
//	err := catService.Create("Moodboard Non Reece Brands", feedName, configName)
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//
//	err = catService.Create("Moodboard Non Reece Categories", feedName, configName)
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//
//	err = catService.Create("TRS Product Color", feedName, configName)
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//
//	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "ff8c36e0-cc3d-48a0-8efe-9a4de800ce14")
//	status, err := contentService.CreateBatch("4c8b4730-7503-485a-9c8e-23af27c61307", feedName, configName)
//	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
//	log.Info(" success created record count  :" + strconv.Itoa(len(status.Success)))
//	if status.FailuresExist() {
//		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
//		status.PrintFailed()
//	}
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//
//	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
//	//err = deleteService.Delete("0e958689-13ce-4eda-9c68-bb4dcc09dd73", "Delete Non Reece Moodboard content", "config_non_reece_product_swatches.yaml")
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//	//
//	//err = deleteService.Delete("0e958689-13ce-4eda-9c68-bb4dcc09dd73", "Delete Non Reece Moodboard images", "config_non_reece_product_swatches.yaml")
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//}
//
func badges(feedName string, configName string) {

	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err := deleteService.Delete("59f1a68b-f518-47a5-83d9-bd16e26c7daa", "Delete Ambulant badge data", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "59f1a68b-f518-47a5-83d9-bd16e26c7daa")
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

func readCustomSearch(feedName string, configName string) {

	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err := deleteService.Delete("59f1a68b-f518-47a5-83d9-bd16e26c7daa", "Delete Ambulant badge data", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "970266a6-c90c-4cc9-8bd6-e5867920e5ce")
	err := contentService.ReadBatch("945ba2f8-5393-41f3-8504-7ea6e5086ed8", feedName, configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

}

func readCustomBadges(feedName string, configName string) {

	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err := deleteService.Delete("59f1a68b-f518-47a5-83d9-bd16e26c7daa", "Delete Ambulant badge data", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "59f1a68b-f518-47a5-83d9-bd16e26c7daa")
	err := contentService.ReadBatch("a86ef75-537b-4c67-aebf-e476e1d2a099", feedName, configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

}

func readTiles(feedName string, configName string) {

	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err := deleteService.Delete("59f1a68b-f518-47a5-83d9-bd16e26c7daa", "Delete Ambulant badge data", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "cbcf8e69-a6ad-474c-b5db-766fc6956cfb")
	err := contentService.ReadBatch("ef48246d-a3e4-4cf7-aee7-32a19d9afd20", feedName, configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

}

func read3dProductConfigurator(feedName string, configName string) {

	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}

	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	//err := deleteService.Delete("59f1a68b-f518-47a5-83d9-bd16e26c7daa", "Delete Ambulant badge data", configName)
	//if err != nil {
	//	errorHandling.WithError(err).Panic(err)
	//}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "cbcf8e69-a6ad-474c-b5db-766fc6956cfb")
	err := contentService.ReadBatch("6972fdb3-f431-4a8b-81f5-95b34148f616", feedName, configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

}

func pig_teams(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "f12876da-0017-41a8-8008-65065ed59f8a")
	status, err := contentService.CreateBatch("06f5cba1-4380-41f1-a115-64b2899e6480", feedName, configName)
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

func pig_projects(feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), "f12876da-0017-41a8-8008-65065ed59f8a")
	status, err := contentService.CreateBatch("ef57950a-dd09-4e3d-924b-0ad5be786091", feedName, configName)
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

//
//func sliBanners(feedName string, configName string) {
//	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
//
//	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
//	//err := deleteService.Delete("24f48506-0498-41f3-8701-43ab4a81b396", "Delete Sli banners", configName)
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//	//err = deleteService.Delete("24f48506-0498-41f3-8701-43ab4a81b396", "Delete Sli banner images", configName)
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//
//	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "24f48506-0498-41f3-8701-43ab4a81b396")
//	status, err := contentService.CreateBatch("92423321-23ce-423e-bff2-74d39a99e449", feedName, configName)
//	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
//	log.Info(" success created record count  :" + strconv.Itoa(len(status.Success)))
//	if status.FailuresExist() {
//		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
//		status.PrintFailed()
//	}
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//}
//
//func brands(feedName string, configName string) {
//	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
//
//	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
//	//err := deleteService.Delete("191c5ef1-9194-4326-bc4f-bc37d681685e", "Delete brands", configName)
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//
//	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "191c5ef1-9194-4326-bc4f-bc37d681685e")
//	status, err := contentService.CreateBatch("9c4cef56-78fc-41d1-b959-d5ddb8f4fd9d", feedName, configName)
//	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
//	log.Info(" success created record count  :" + strconv.Itoa(len(status.Success)))
//	if status.FailuresExist() {
//		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
//		status.PrintFailed()
//	}
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//}
//
//func customSearch(feedName string, configName string) {
//	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
//
//	//deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
//	//err := deleteService.Delete("191c5ef1-9194-4326-bc4f-bc37d681685e", "Delete brands", configName)
//	//if err != nil {
//	//	errorHandling.WithError(err).Panic(err)
//	//}
//
//	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "970266a6-c90c-4cc9-8bd6-e5867920e5ce")
//	status, err := contentService.CreateBatch("9c95a805-99a2-41a4-b423-0bdfa592cc53", feedName, configName)
//	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
//	log.Info(" success created record count  :" + strconv.Itoa(len(status.Success)))
//	if status.FailuresExist() {
//		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
//		status.PrintFailed()
//	}
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//}
//
//func createBrands() {
//	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
//	var err error
//	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
//	err = catService.Create("Reece AU website product brand", "brands.csv", "config_reece_brands.yaml")
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//}
//
func createCats(catName string, feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	err = catService.Create(catName, feedName, configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}
}

//
//func deleteCats() {
//	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
//	var err error
//	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
//	err = catService.Delete("Reece AU website category")
//	if err != nil {
//		errorHandling.WithError(err).Panic(err)
//	}
//}

func main() {
	//feedName := "20221205_materials.csv"
	//configName := "config_3d_product_material.yaml"
	////delete_reece(configName)
	//threed_materials(feedName, configName)

	//feedName := "20221205_Metadata_acoustic.csv"
	//configName := "config_3d_product_configurator.yaml"
	////delete_reece(configName)
	//threed_materials_configuration(feedName, configName)

	//feedName := "20220718_Metadata_acoustic.csv"
	//configName := "config_3d_product_configurator_update.yaml"
	////delete_reece(configName)
	//threed_materials_configuration_update(feedName, configName)

	//feedName := "20220607_metadata.csv"
	//configName := "config_3d_product_material_update.yaml"
	////delete_reece(configName)
	//threed_materials_update(feedName, configName)

	//feedName := "20220808_Metadata_acoustic.csv"
	//configName := "config_3d_product_configurator.yaml"
	////delete_reece(configName)
	//threed_materials_configuration(feedName, configName)

	//feedName := "Data for URL codes.csv"
	//configName := "custom_search_all.yaml"
	////delete_reece(configName)
	//custom_search(feedName, configName)

	//feedName := "Data for URL codes.csv"
	//configName := "custom_search_products_all.yaml"
	////delete_reece(configName)
	//custom_search_products(feedName, configName)

	//feedName := "20221104_custom_order_nz_remove.csv"
	//configName := "badge_custom_order_products_delete.yaml"
	////delete_reece(configName)
	//badges(feedName, configName)

	//feedName := "20221205_materials.csv"
	//configName := "config_3d_product_material_review.yaml"
	////delete_reece(configName)
	//threed_materials_review(feedName, configName)

	//feedName := "20221205_Metadata_acoustic.csv"
	//configName := "config_3d_product_configurator_review.yaml"
	////delete_reece(configName)
	//threed_materials_configuration_review(feedName, configName)

	//feedName := "National Tiles - 3d Planner & Moodboard - Export to Accoustic.csv"
	//configName := "3d_planner_tiles_au_insert.yaml"
	////delete_reece(configName)
	//threed_tiles_update(feedName, configName)

	//feedName := "tiles_in_acoustic_20220829.csv"
	//configName := "3d_planner_tiles_nz_insert.yaml"
	////delete_reece(configName)
	//threed_tiles_update(feedName, configName)

	//feedName := "new_arrivals_to_delete.csv"
	//configName := "custom_search_new_arrivals_read.yaml"
	////delete_reece(configName)
	//readCustomSearch(feedName, configName)

	//feedName := "custom_order_badges_to_delete.csv"
	//configName := "custom_order_to_delete_read.yaml"
	////delete_reece(configName)
	//readCustomBadges(feedName, configName)

	//feedName := "tiles_in_acoustic_20220829.csv"
	//configName := "3d_planner_tiles_read.yaml"
	////delete_reece(configName)
	//readTiles(feedName, configName)

	//feedName := "3d_product_configurations_to_update.csv"
	//configName := "config_3d_product_configurator_read.yaml"
	////delete_reece(configName)
	//read3dProductConfigurator(feedName, configName)

	feedName := "pig-teams.csv"
	configName := "pig-teams.yaml"
	//delete_pig_all(configName)
	pig_teams(feedName, configName)

	//feedName := "Imagin3D Paint Colours_20221108.csv"
	//configName := "3d_planner_paint_insert.yaml"
	//createCats("3D Planner Paint Category", feedName, configName)

	//feedName := "Imagin3D Paint Colours_20221108.csv"
	//configName := "3d_planner_paint_insert.yaml"
	////delete_reece(configName)
	//threed_paint(feedName, configName)

}
