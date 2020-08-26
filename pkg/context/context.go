package context

import "os"

func IsDebugEnabled () bool {
	return  os.Getenv("DebugEnabled") == "true"
}

func AcousticAuthUrl () string {
	return os.Getenv("AcousticAuthURL")
}

func AcousticAPIUrl () string {
	return os.Getenv("AcousticAPIURL")
}

func AcousticAuthUserName () string {
	return os.Getenv("AcousticAuthUserName")
}

func AcousticAuthPassword () string {
	return os.Getenv("AcousticAuthPassword")
}
