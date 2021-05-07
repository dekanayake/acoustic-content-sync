package cmd

import "github.com/spf13/cobra"
import log "github.com/sirupsen/logrus"

func getFlagStringValue(cmd *cobra.Command, flagName string) string {
	value, err := cmd.Flags().GetString(flagName)
	if err != nil {
		log.Panic("Error occurred while getting the flag value :"+flagName, err)
	}
	return value
}

func getFlagBoolValue(cmd *cobra.Command, flagName string) bool {
	value, err := cmd.Flags().GetBool(flagName)
	if err != nil {
		log.Panic("Error occurred while getting the flag value :"+flagName, err)
	}
	return value
}
