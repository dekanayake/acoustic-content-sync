package cmd

import (
	"errors"
	"fmt"
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/csv"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/logrus"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var publishContent bool

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create acoustic content",
	Long:  `Create acoustic content with assets`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a content type argument")
		}
		err := validateContentType(cmd, args[0])
		if err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Create aciustic content")
	},
}

func validateContentType(cmd *cobra.Command, contentType string) error {
	configLocation := getFlagStringValue(cmd, "AcousticConfigLocation")
	config, err := csv.InitConfig(configLocation)
	if err != nil {
		return err
	}
	configTypeMapping, _ := config.GetContentType(contentType)
	if configTypeMapping == nil {
		return errors.New(" content type is not available in config : " + contentType)
	}
	return nil
}

func executeCreateCommand(cmd *cobra.Command, args []string) error {
	publishContent := getFlagBoolValue(cmd, "PublishContent")
	if publishContent {
		os.Setenv("ContentStatus", "ready")
	} else {
		os.Setenv("ContentStatus", "draft")
	}
	writeErrorsToFile := getFlagBoolValue(cmd, "WriteErrorsToFile")
	if writeErrorsToFile {
		os.Setenv("ErrorLogFileLocation", "error.log")
	}
	contentType := args[0]
	feedLocation := getFlagStringValue(cmd, "DataCSVFileLocation")
	configFileLocation := getFlagStringValue(cmd, "AcousticConfigLocation")
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	contentService := csv.NewContentService(os.Getenv("AcousticAPIURL"), "ff8c36e0-cc3d-48a0-8efe-9a4de800ce14")
	status, err := contentService.Create(contentType, feedLocation, configFileLocation)
	log.Info(" total records :" + strconv.Itoa(status.TotalCount()))
	log.Info(" success created record count  :" + strconv.Itoa(len(status.Success)))
	if status.FailuresExist() {
		log.Error("There are " + strconv.Itoa(len(status.Failed)) + " failures in creating contents , please check the log in " + env.ErrorLogFileLocation())
		status.PrintFailed()
	}
	if err != nil {
		errorHandling.WithError(err)
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.LocalFlags().BoolVar(&publishContent, "PublishContent", true, "Publish content after creation complete, false will create the content in draft mode")
}
func init() {

}
