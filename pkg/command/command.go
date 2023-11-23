package command

import (
	"flag"
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/csv"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/thoas/go-funk"
	"os"
	"strconv"
	"strings"
	"time"
)
import log "github.com/sirupsen/logrus"

type Argument struct {
	name      string
	usage     string
	mandatory bool
}

type CommandHandler func(params map[string]*string) error

type Command struct {
	name      string
	arguments []Argument
	handler   CommandHandler
}

type Action struct {
	command        Command
	argumentValues map[string]string
}

func (command *Command) execute() error {
	argumentValues := make(map[string]*string)
	arguments := command.arguments
	for _, argument := range arguments {
		argValue := flag.String(argument.name, "", argument.usage)
		flag.Parse()
		if argument.mandatory && *argValue == "" {
			log.Error("Please provide the " + argument.usage + ", usage: -" + argument.name + "=...")
		}
		log.Info(argument.usage + " : " + *argValue)
		argumentValues[argument.name] = argValue
	}
	return command.handler(argumentValues)
}

type ActionProcessor struct {
	commands []Command
}

func (processor *ActionProcessor) Execute() error {
	contentOperation := flag.String("operation", "", "Operation")
	flag.Parse()
	if len(strings.TrimSpace(*contentOperation)) == 0 {
		operations := funk.Reduce(processor.commands, func(acc string, cmd string) string {
			return acc + cmd
		}, "").(string)
		log.Error("Please provide the Content Operation. Commands supported : " + operations)
		os.Exit(1)
	}

	matchedContentOperationCmd := funk.Filter(processor.commands, func(command Command) bool {
		return command.name == *contentOperation
	})

	if matchedContentOperationCmd == "" {
		log.Error("Matching command not found for command :" + *contentOperation)
		os.Exit(1)
	}

	command := matchedContentOperationCmd.(Command)
	return command.execute()

}

func createOrUpdateContents(feedName string, configName string, acousticContentLib string, contentType string) error {
	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), acousticContentLib)
	status, err := contentService.CreateBatch(contentType, feedName, configName)
	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
	log.Info(" success created record count  :" + strconv.Itoa(len(status.Success)))
	if status.FailuresExist() {
		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
		status.PrintFailed()
	}
	return err
}

func deleteContents(deleteUsingFeed bool, deleteMappingName string, feedName string, configName string, libraryID string, contentType string) error {
	deleteService := csv.NewDeleteService(env.AcousticAPIUrl())
	if deleteUsingFeed {
		status, err := deleteService.DeleteByFeed(deleteMappingName, contentType, feedName, configName)
		log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
		log.Info(" success deleted record count  :" + strconv.Itoa(len(status.Success)))
		if status.FailuresExist() {
			log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
			status.PrintFailed()
		}
		return err
	} else {
		err := deleteService.Delete(libraryID, deleteMappingName, configName)
		return err
	}
}

func createCategories(catName string, feedName string, configName string) error {
	catService := csv.NewCategoryService(os.Getenv("AcousticAPIURL"))
	err := catService.Create(catName, feedName, configName)
	return err
}

func createSitePages(siteId string, parentPageId string, contentType string, dataFeedPath string, configPath string) error {
	os.Setenv("ParentPageContentTypeID", contentType)
	siteUseCase := csv.NewSiteUseCase(os.Getenv("AcousticAPIURL"))
	status, err := siteUseCase.CreatePages(siteId, parentPageId, contentType, dataFeedPath, configPath)
	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
	log.Info(" success created pages count  :" + strconv.Itoa(len(status.Success)))
	if status.FailuresExist() {
		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating pages , please check the log in " + env.ErrorLogFileLocation())
		status.PrintFailed()
	}
	return err
}

func createPageForContent(siteId string, parentPageId string, contentID string, contentType string, relativeUrl string) error {
	os.Setenv("ParentPageContentTypeID", contentType)
	siteUseCase := csv.NewSiteUseCase(os.Getenv("AcousticAPIURL"))
	createdPageID, err := siteUseCase.CreatePageForContent(siteId, parentPageId, contentID, relativeUrl)
	log.Info("Page created with ID :" + createdPageID)
	return err
}

func clone(id string, sourceAcousticAuthApiHost string, sourceAcousticAPIKey string, targetAcousticAuthApiHost string, targetAcousticAPIKey string) error {
	copyUseCase := csv.NewContentCopyUserCase(sourceAcousticAuthApiHost, sourceAcousticAPIKey, targetAcousticAuthApiHost, targetAcousticAPIKey)

	_, err := copyUseCase.CopyContent(id, "_CL:"+time.Now().Format(time.ANSIC))
	return err
}

func readContents(feedName string, configName string, acousticContentLib string, contentType string) error {
	contentService := csv.NewContentUseCase(os.Getenv("AcousticAPIURL"), acousticContentLib)
	err := contentService.ReadBatch(contentType, feedName, configName)
	return err
}

func NewCommandActionProcessor() *ActionProcessor {
	createCommand := Command{
		name: "CREATE",
		arguments: []Argument{
			Argument{
				name:      "feedLocation",
				usage:     "File path of the feed",
				mandatory: true,
			},
			Argument{
				name:      "configLocation",
				usage:     "File path of the config",
				mandatory: true,
			},
			Argument{
				name:      "acousticLibraryID",
				usage:     "Acoustic Library ID",
				mandatory: true,
			},
			Argument{
				name:      "contentTypeID",
				usage:     "Content Type ID",
				mandatory: true,
			},
		},
		handler: func(params map[string]*string) error {
			return createOrUpdateContents(*params["feedLocation"], *params["configLocation"], *params["acousticLibraryID"], *params["contentTypeID"])
		},
	}

	updateCommand := Command{
		name: "UPDATE",
		arguments: []Argument{
			Argument{
				name:      "feedLocation",
				usage:     "File path of the feed",
				mandatory: true,
			},
			Argument{
				name:      "configLocation",
				usage:     "File path of the config",
				mandatory: true,
			},
			Argument{
				name:      "acousticLibraryID",
				usage:     "Acoustic Library ID",
				mandatory: true,
			},
			Argument{
				name:      "contentTypeID",
				usage:     "Content Type ID",
				mandatory: true,
			},
		},
		handler: func(params map[string]*string) error {
			return createOrUpdateContents(*params["feedLocation"], *params["configLocation"], *params["acousticLibraryID"], *params["contentTypeID"])
		},
	}

	readCommand := Command{
		name: "READ",
		arguments: []Argument{
			Argument{
				name:      "feedLocation",
				usage:     "File path of the feed",
				mandatory: true,
			},
			Argument{
				name:      "configLocation",
				usage:     "File path of the config",
				mandatory: true,
			},
			Argument{
				name:      "acousticLibraryID",
				usage:     "Acoustic Library ID",
				mandatory: true,
			},
			Argument{
				name:      "contentTypeID",
				usage:     "Content Type ID",
				mandatory: true,
			},
		},
		handler: func(params map[string]*string) error {
			return readContents(*params["feedLocation"], *params["configLocation"], *params["acousticLibraryID"], *params["contentTypeID"])
		},
	}

	createCategoryCommand := Command{
		name: "CREATE_CATEGORY",
		arguments: []Argument{
			Argument{
				name:      "categoryName",
				usage:     "Root category name to add new categories",
				mandatory: true,
			},
			Argument{
				name:      "feedLocation",
				usage:     "File path of the feed",
				mandatory: true,
			},
			Argument{
				name:      "configLocation",
				usage:     "File path of the config",
				mandatory: true,
			},
		},
		handler: func(params map[string]*string) error {
			return createCategories(*params["categoryName"], *params["feedLocation"], *params["configLocation"])
		},
	}

	createSitePages := Command{
		name: "CREATE_SITE_PAGES",
		arguments: []Argument{
			Argument{
				name:      "siteId",
				usage:     "Site ID",
				mandatory: true,
			},
			Argument{
				name:      "parentPageID",
				usage:     "Parent page ID",
				mandatory: true,
			},
			Argument{
				name:      "contentTypeID",
				usage:     "Content Type ID",
				mandatory: true,
			},
			Argument{
				name:      "feedLocation",
				usage:     "File path of the feed",
				mandatory: true,
			},
			Argument{
				name:      "configLocation",
				usage:     "File path of the config",
				mandatory: true,
			},
		},
		handler: func(params map[string]*string) error {
			return createSitePages(*params["siteId"], *params["parentPageID"], *params["contentTypeID"], *params["feedLocation"], *params["configLocation"])
		},
	}

	createSitePageForContent := Command{
		name: "CREATE_SITE_PAGE_FOR_CONTENT",
		arguments: []Argument{
			Argument{
				name:      "siteId",
				usage:     "Site ID",
				mandatory: true,
			},
			Argument{
				name:      "parentPageID",
				usage:     "Parent page ID",
				mandatory: true,
			},
			Argument{
				name:      "contentIDForPage",
				usage:     "Content ID to create page",
				mandatory: true,
			},
			Argument{
				name:      "contentTypeID",
				usage:     "Content Type ID",
				mandatory: true,
			},
			Argument{
				name:      "relativeUrlOfPage",
				usage:     "Relative URL of the page",
				mandatory: true,
			},
		},
		handler: func(params map[string]*string) error {
			return createPageForContent(*params["siteId"], *params["parentPageID"], *params["contentIDForPage"], *params["contentTypeID"], *params["relativeUrlOfPage"])
		},
	}

	cloneContent := Command{
		name: "CLONE_CONTENT",
		arguments: []Argument{
			Argument{
				name:      "idToClone",
				usage:     "ID to clone",
				mandatory: true,
			},
			Argument{
				name:      "sourceAcousticAuthAPIHost",
				usage:     "Source Acoustic  Host API",
				mandatory: true,
			},
			Argument{
				name:      "sourceAcousticAPIKey",
				usage:     "Source Acoustic API Key",
				mandatory: true,
			},
			Argument{
				name:      "targetAcousticAuthAPIHost",
				usage:     "Target Acoustic  Host API",
				mandatory: true,
			},
			Argument{
				name:      "targetAcousticAPIKey",
				usage:     "Target Acoustic API Key",
				mandatory: true,
			},
		},
		handler: func(params map[string]*string) error {
			return clone(*params["idToClone"], *params["sourceAcousticAuthAPIHost"], *params["sourceAcousticAPIKey"], *params["targetAcousticAuthAPIHost"], *params["targetAcousticAPIKey"])
		},
	}

	return &ActionProcessor{
		commands: []Command{createCommand, updateCommand, readCommand, createCategoryCommand, createSitePages, createSitePageForContent, cloneContent},
	}
}
