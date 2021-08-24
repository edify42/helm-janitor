package config

import (
	"os"
)

var (
	AppName       = "helm-janitor"
	TTLKey        = "janitor/ttl"
	ExpiryKey     = "janitor/expires"
	AnnotationKey = "janitorAnnotations"
	DefaultTTL    = 7 * 24 * 60 * 60 // 7 days in seconds.
)

// GetenvWithDefault
func GetenvWithDefault(env string, def string) string {
	if os.Getenv(env) != "" {
		return os.Getenv(env)
	}
	return def
}

// GetenvWithDefaultBool
func GetenvWithDefaultBool(env string, def bool) bool {
	a := os.Getenv(env)
	if a != "" {
		matchedTrue := []string{"TRUE", "true", "1", "true"}
		for _, b := range matchedTrue {
			if a == b {
				return true
			}
		}
		return false
	}
	return def
}
