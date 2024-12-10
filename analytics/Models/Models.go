package models

type KafkaConfig struct {
	Host  string
	Port  string
	Topic string
}

type AnalyticsData struct {
	IPAddress string `json:"ip_address"`
	Referrer  string `json:"referrer"`
	UserAgent string `json:"user_agent"`
	ShortURL  string `json:"short_url"`
	Timestamp string `json:"timestamp"`
}

type ProcessedData struct {
	Country   string `json:"country,omitempty"`
	City      string `json:"city,omitempty"`
	OS        string `json:"os,omitempty"`
	Device    string `json:"device,omitempty"`
	Browser   string `json:"browser,omitempty"`
	Timestamp string `json:"timestamp"`
	Referrer  string `json:"referrer,omitempty"`
	ShortURL  string `json:"short_url"`
}

type IP struct {
	Country string `json:"country,omitempty"`
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
}
