package main

import (
	"flag"
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/csv"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/logrus"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"strings"
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

func createOrUpdateContents(feedName string, configName string, acousticContentLib string, contentType string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	var err error
	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), acousticContentLib)
	status, err := contentService.CreateBatch(contentType, feedName, configName)
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

func deleteContents(deleteUsingFeed bool, deleteMappingName string, feedName string, configName string, libraryID string, contentType string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	if deleteUsingFeed {
		status, err := deleteService.DeleteByFeed(deleteMappingName, contentType, feedName, configName)
		if err != nil {
			errorHandling.WithError(err).Panic(err)
		}
		log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
		log.Info(" success deleted record count  :" + strconv.Itoa(len(status.Success)))
		if status.FailuresExist() {
			log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
			status.PrintFailed()
		}
		if err != nil {
			errorHandling.WithError(err).Panic(err)
		}
	} else {
		err := deleteService.Delete(libraryID, deleteMappingName, configName)
		if err != nil {
			errorHandling.WithError(err).Panic(err)
		}
	}
}

func createCategories(catName string, feedName string, configName string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	err := catService.Create(catName, feedName, configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}
}

func createSitePages(siteId string, parentPageId string, contentType string, dataFeedPath string, configPath string) {
	os.Setenv("ParentPageContentTypeID", contentType)
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	siteUseCase := csv.NewSiteUseCase(os.Getenv("AcousticAPIURL"))
	status, err := siteUseCase.CreatePages(siteId, parentPageId, contentType, dataFeedPath, configPath)
	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
	log.Info(" success created pages count  :" + strconv.Itoa(len(status.Success)))
	if status.FailuresExist() {
		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating pages , please check the log in " + env.ErrorLogFileLocation())
		status.PrintFailed()
	}
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}
}

func createPageForContent(siteId string, parentPageId string, contentID string, relativeUrl string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	siteUseCase := csv.NewSiteUseCase(os.Getenv("AcousticAPIURL"))
	createdPageID, err := siteUseCase.CreatePageForContent(siteId, parentPageId, contentID, relativeUrl)
	log.Info("Page created with ID :" + createdPageID)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}
}

func clone(id string, libraryID string) {
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	copyUseCase := csv.NewContentCopyUserCase(os.Getenv("AcousticAPIURL"))
	_, err := copyUseCase.CopyContent(id, libraryID)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}
	//log.Info(" total records :" + strconv.Itoa(contentStatus.TotalCount()))
}

func readContents(feedName string, configName string, acousticContentLib string, contentType string) {

	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}

	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), acousticContentLib)
	err := contentService.ReadBatch(contentType, feedName, configName)
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}

}

func execute() {
	log.Info("--------------Running Synky CLI----------------")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	feedLocation := flag.String("feedLocation", "", "File path of the feed")
	configLocation := flag.String("configLocation", "", "File path of the config")
	acousticLibraryID := flag.String("acousticLibraryID", "", "Acoustic Library ID")
	contentTypeID := flag.String("contentTypeID", "", "Content Type ID")
	contentOperation := flag.String("operation", "", "Operation")
	deleteByFeed := flag.Bool("deleteByFeed", false, "Delete By Feed")
	categoryName := flag.String("categoryName", "", "Root category name to add new categories")
	deleteMappingName := flag.String("deleteMappingName", "", "Delete Mapping Name")
	siteId := flag.String("siteID", "", "Site ID")
	parentPageID := flag.String("parentPageID", "", "Parent page ID")
	idToClone := flag.String("idToClone", "", "ID to clone")
	contentIDForPage := flag.String("contentIDForPage", "", "Content ID to create page")
	relativeUrlOfPage := flag.String("relativeUrlOfPage", "", "Relative URL of the page")
	flag.Parse()

	log.Info("feed location :" + *feedLocation)
	log.Info("config location :" + *configLocation)
	log.Info("Acoustic Library ID :" + *acousticLibraryID)
	log.Info("Content type ID :" + *contentTypeID)
	log.Info("Root Category Name :" + *categoryName)
	log.Info("Operation :" + *contentOperation)
	log.Info("Site ID :" + *siteId)
	log.Info("Parent Page ID :" + *parentPageID)
	log.Info("ID to clone :" + *idToClone)
	log.Info("Content ID to create page :" + *contentIDForPage)
	log.Info("Relative URL of the page :" + *relativeUrlOfPage)

	if len(strings.TrimSpace(*contentOperation)) == 0 {
		log.Error("Please provide the Content Operation (CREATE for create , UPDATE for update , READ for read) ")
		os.Exit(1)
	}

	if len(strings.TrimSpace(*feedLocation)) == 0 && *contentOperation != "CLONE_CONTENT" && *contentOperation != "CREATE_SITE_PAGE_FOR_CONTENT" {
		log.Error("Please provide the feed location")
		os.Exit(1)
	}

	if len(strings.TrimSpace(*configLocation)) == 0 && *contentOperation != "CLONE_CONTENT" && *contentOperation != "CREATE_SITE_PAGE_FOR_CONTENT" {
		log.Error("Please provide the config location")
		os.Exit(1)
	}

	if len(strings.TrimSpace(*acousticLibraryID)) == 0 && *contentOperation != "CREATE_CATEGORY" && *contentOperation != "CREATE_SITE_PAGE_FOR_CONTENT" {
		log.Error("Please provide the Acoustic Library ID")
		os.Exit(1)
	} else {
		os.Setenv("LibraryID", strings.TrimSpace(*acousticLibraryID))
	}

	if len(strings.TrimSpace(*contentTypeID)) == 0 && *contentOperation != "CREATE_CATEGORY" && *contentOperation != "CLONE_CONTENT" && *contentOperation != "CREATE_SITE_PAGE_FOR_CONTENT" {
		log.Error("Please provide the Content Type ID")
		os.Exit(1)
	}

	if len(strings.TrimSpace(*idToClone)) == 0 && *contentOperation == "CLONE_CONTENT" {
		log.Error("Please provide the Content ID")
		os.Exit(1)
	}

	if *contentOperation == "CREATE" || *contentOperation == "UPDATE" {
		createOrUpdateContents(*feedLocation, *configLocation, *acousticLibraryID, *contentTypeID)
	} else if *contentOperation == "READ" {
		readContents(*feedLocation, *configLocation, *acousticLibraryID, *contentTypeID)
	} else if *contentOperation == "DELETE" {
		deleteContents(*deleteByFeed, *deleteMappingName, *feedLocation, *configLocation, *acousticLibraryID, *contentTypeID)
	} else if *contentOperation == "CREATE_CATEGORY" {
		createCategories(*categoryName, *feedLocation, *configLocation)
	} else if *contentOperation == "CREATE_SITE_PAGES" {
		createSitePages(*siteId, *parentPageID, *contentTypeID, *feedLocation, *configLocation)
	} else if *contentOperation == "CREATE_SITE_PAGE_FOR_CONTENT" {
		createPageForContent(*siteId, *parentPageID, *contentIDForPage, *relativeUrlOfPage)
	} else if *contentOperation == "CLONE_CONTENT" {
		clone(*idToClone, *acousticLibraryID)
	} else {
		log.Error("Please provide the Content Operation (CREATE for create , UPDATE for update , READ for read , provided operation : {}", *contentOperation)
		os.Exit(1)
	}

}

func main() {
	execute()
}
