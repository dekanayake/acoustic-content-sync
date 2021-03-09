package env

import (
	"log"
	"os"
)

func GetOrPanic(variable string) string {
	varValue := os.Getenv(variable)
	if varValue == "" {
		log.Panic("Env variable not available :" + variable)
	}
	return varValue
}

func Get(variable string) string {
	varValue := os.Getenv(variable)
	return varValue
}

func IsDebugEnabled() bool {
	return os.Getenv("DebugEnabled") == "true"
}

func AcousticAuthUrl() string {
	return GetOrPanic("AcousticAuthURL")
}

func AcousticAPIUrl() string {
	return GetOrPanic("AcousticAPIURL")
}

func LibraryID() string {
	return GetOrPanic("LibraryID")
}

func AcousticAuthUserName() string {
	return Get("AcousticAuthUserName")
}

func AcousticAuthPassword() string {
	return Get("AcousticAuthPassword")
}

func AcousticAPIKey() string {
	return Get("AcousticAPIKey")
}

func ContentStatus() string {
	return GetOrPanic("ContentStatus")
}

func CategoryHierarchySeperator() string {
	return GetOrPanic("CategoryHierarchySeperator")
}

func MultipleItemsSeperator() string {
	return GetOrPanic("MultipleItemsSeperator")
}

func LogErrorsToFile() bool {
	return GetOrPanic("WriteErrorsToFile") == "true"
}

func ErrorLogFileLocation() string {
	return GetOrPanic("ErrorLogFileLocation")
}
