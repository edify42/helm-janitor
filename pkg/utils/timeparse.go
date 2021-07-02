package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

// Time is
type Time struct {
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

// ParseTime expects one of the following:
// '2d' = 2 days; '2h' 2 hours; '2m' = 2 minutes; '2s' = 2 seconds
// Number followed by d, h, m or s.
// Returned struct Time as an approximation of the input.
func ParseTime(input string) (Time, error) {
	re := regexp.MustCompile(`([0-9]+)([dhms])`)
	rtn := Time{}
	if re.MatchString(input) {
		res := re.FindStringSubmatch(input)
		value, _ := strconv.Atoi(res[1])
		unit := res[2]
		if unit == "d" {
			rtn.Days = value
			rtn.Hours = rtn.Days * 24
			rtn.Minutes = rtn.Hours * 60
			rtn.Seconds = rtn.Minutes * 60
		} else if unit == "h" {
			rtn.Days = value / 24
			rtn.Hours = value
			rtn.Minutes = value * 60
			rtn.Seconds = rtn.Minutes * 60
		} else if unit == "m" {
			rtn.Days = value / 60 / 24
			rtn.Hours = value / 60
			rtn.Minutes = value
			rtn.Seconds = value * 60
		} else {
			rtn.Days = value / 60 / 60 / 24
			rtn.Hours = value / 60 / 60
			rtn.Minutes = value / 60
			rtn.Seconds = value
		}
	}
	return rtn, fmt.Errorf("input %s not valid", input)
}
