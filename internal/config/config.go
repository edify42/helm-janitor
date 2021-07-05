package config

var (
	AppName       = "helm-janitor"
	TTLKey        = "janitor/ttl"
	ExpiryKey     = "janitor/expiry"
	AnnotationKey = "janitorAnnotations"
	DefaultTTL    = 7 * 24 * 60 * 60 // 7 days in seconds.
)
