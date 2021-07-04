package scan

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// ValidateScanArgs returns true when we have inputs
func ValidateScanArg(args []string) bool {
	if len(args) == 0 {
		return false
	}
	if len(args) > 1 {
		log.Fatalf("too many arguments provided: %v", args)
	}

	a := strings.Split(args[0], ",")

	for _, v := range a {
		i := strings.Split(v, "=")
		if len(i) != 2 {
			log.Fatalf("incorrect 'key=value' selector usage: %s. fix input %s", v, args)
		}
	}
	return true
}
