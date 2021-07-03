package config

var (
	AppName    = "helm-janitor"
	TTLKey     = "helm-janitor/ttl"
	ExpiryKey  = "helm-janitor/expiry"
	DefaultTTL = 7 * 24 * 60 * 60 // 7 days in seconds.
)
