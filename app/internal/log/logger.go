package internal_log

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "", 0)
