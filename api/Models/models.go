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

type AnalyticsData struct {
	IPAddress string `json:"ip_address"`
	Referrer  string `json:"referrer"`
	UserAgent string `json:"user_agent"`
	ShortURL  string `json:"short_url"`
	Timestamp string `json:"timestamp"`
}

func NewAnalyticsData(ip, referrer, userAgent, shortURL string) AnalyticsData {
	return AnalyticsData{
		IPAddress: defaultIfEmpty(ip, "Not Available"),
		Referrer:  defaultIfEmpty(referrer, "Not Available"),
		UserAgent: defaultIfEmpty(userAgent, "Not Available"),
		ShortURL:  defaultIfEmpty(shortURL, "Not Available"),
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func defaultIfEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
