package lambda

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// runs the generic handler to execute helm delete...
// when the ttl expires.

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	logLevel := "info"
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Errorf("Dodgy log level set: %s", logLevel)
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(level)
	}
}

func main() {
	log.Infof("starting")
}
