package cmd

import (
	"errors"
	"fmt"
	"github.com/bgentry/speakeasy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var acousticAuthURL string
var acousticAPIURL string
var libraryID string
var authenticateViaAPIKey bool
var acousticAPIKey string
var acousticAuthUserName string
var writeErrorsToFile bool
var dataCsvFileLocation string
var acousticConfigLocation string

var rootCmd = &cobra.Command{
	Use:   "acoustic-content-sync",
	Short: "acoustic-content-sync helps you to ingest contents in CSV to Acoustic Content Headless CMS",
	Long:  `Ingest your contents in CSV to Acoustic Content Headless CMS Supports content/asset creation, deletion, taxonomy creation`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeCommand(cmd, args)
	},
}

func preRunValidation(cmd *cobra.Command, args []string) error {
	configFileLocation := getFlagStringValue(cmd, "AcousticConfigLocation")
	_, err := os.Open(configFileLocation)
	if err != nil {
		return errors.New("Failed in opening the config file : " + configFileLocation + " error :" + err.Error())
	}
	dataCsvFileLocation := getFlagStringValue(cmd, "DataCSVFileLocation")
	_, err = os.Open(dataCsvFileLocation)
	if err != nil {
		return errors.New("Failed in opening the csv file : " + dataCsvFileLocation + " error :" + err.Error())
	}
	return nil

}

func setAuthentication(cmd *cobra.Command, args []string) error {
	authViaAPIKey := getFlagBoolValue(cmd, "AuthenticateViaAPIKey")
	if authViaAPIKey {
		apiKey := getFlagStringValue(cmd, "AcousticAPIKey")
		if apiKey == "" {
			return errors.New("No Acoustic Key provided, please provide the acoustic API")
		}
		os.Setenv("AcousticAPIKey", apiKey)
	} else {
		authUserName := getFlagStringValue(cmd, "AcousticAuthUserName")
		if authUserName == "" {
			return errors.New("No Acoustic user name provided, please provide the acoustic auth user name")
		} else {
			prompt1 := fmt.Sprintf("Password of %s: ", authUserName)
			password, _ := speakeasy.Ask(prompt1)
			if password == "" {
				return errors.New("Please provide the password")
			}
			os.Setenv("AcousticAuthUserName", authUserName)
			os.Setenv("AcousticAuthPassword", password)
		}
	}
	return nil
}

func executeCommand(cmd *cobra.Command, args []string) error {
	err := preRunValidation(cmd, args)
	if err != nil {
		return err
	}
	err = setAuthentication(cmd, args)
	if err != nil {
		return err
	}
	setStringEnvVariable(cmd, "AcousticAuthURL", "AcousticAuthURL")
	setStringEnvVariable(cmd, "AcousticAuthURL", "AcousticAuthURL")
	setBoolEnvVariable(cmd, "WriteErrorsToFile", "WriteErrorsToFile")
	log.Info("Running")
	return nil
}

func setStringEnvVariable(cmd *cobra.Command, envVarName string, flagName string) {
	os.Setenv(envVarName, getFlagStringValue(cmd, flagName))
}

func setBoolEnvVariable(cmd *cobra.Command, envVarName string, flagName string) {
	os.Setenv(envVarName, strconv.FormatBool(getFlagBoolValue(cmd, flagName)))
}

func init() {
	rootCmd.PersistentFlags().StringVar(&acousticAuthURL, "AcousticAuthURL", "", "Acoustic Auth URL")
	rootCmd.MarkPersistentFlagRequired("AcousticAuthURL")
	rootCmd.PersistentFlags().StringVar(&acousticAPIURL, "AcousticAPIURL", "", "Acoustic API URL")
	rootCmd.MarkPersistentFlagRequired("AcousticAPIURL")
	rootCmd.PersistentFlags().BoolVar(&authenticateViaAPIKey, "AuthenticateViaAPIKey", true, "Authenticate via APIKey")
	rootCmd.PersistentFlags().StringVar(&acousticAPIKey, "AcousticAPIKey", "", "Acoustic API Key")
	rootCmd.PersistentFlags().StringVar(&acousticAuthUserName, "AcousticAuthUserName", "", "Acoustic Auth user name")
	rootCmd.PersistentFlags().BoolVar(&writeErrorsToFile, "WriteErrorsToFile", true, "Write Errors to log file")
	rootCmd.PersistentFlags().StringVar(&dataCsvFileLocation, "DataCSVFileLocation", "", "Data CSV File location")
	rootCmd.MarkPersistentFlagRequired("DataCSVFileLocation")
	rootCmd.PersistentFlags().StringVar(&acousticConfigLocation, "AcousticConfigLocation", "", "Acoustic config yaml location")
	rootCmd.MarkPersistentFlagRequired("AcousticConfigLocation")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
