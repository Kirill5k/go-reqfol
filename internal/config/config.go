package config

type Client struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeoutMs   int
	TimeoutMs           int
	RetryCount          int
	RetryWaitTimeMs     int
	RetryMaxWaitTimeMs  int
}
