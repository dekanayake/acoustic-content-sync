package context

import "os"

func IsDebugEnabled () bool {
	return  os.Getenv("DebugEnabled") == "true"
}
