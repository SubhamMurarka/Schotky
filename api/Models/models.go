package models

import "time"

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	URL             string        `json:"url"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}
