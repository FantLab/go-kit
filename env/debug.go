package env

import "os"

var debug = os.Getenv("DEBUG") != ""

func IsDebug() bool {
	return debug
}
