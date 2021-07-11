package scan

import (
	"fmt"
	"strings"
)

// ValidateScanArgs returns true when we have inputs
func ValidateScanArg(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	if len(args) > 1 {
		return false, fmt.Errorf("too many arguments provided: %v", args)
	}

	a := strings.Split(args[0], ",") // e.g. `key1=value1,key2=value2`

	for _, v := range a {
		i := strings.Split(v, "=")
		if len(i) != 2 {
			return false, fmt.Errorf("incorrect 'key=value' selector usage: %s. fix input %s", v, args) // e.g. key=value1=value2 fails
		}
	}
	return true, nil
}
