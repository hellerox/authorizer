package logfile

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Init creates a file named authorizer.log and log all the application errors in there
func Init() {
	log.SetLevel(log.DebugLevel)

	f, err := os.OpenFile("authorizer.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		log.Errorln("error writing file: ", err)
	}

	log.SetOutput(f)

	log.Println("---------------")
}
