package models

import "time"

type Site struct {
	ID int `json:"id"`
	URL string `json:"url"`
	IntervalSeconds int `json:"interval_seconds"`
	CreatedAt time.Time `json:"created_at"`
}

type CheckLog struct {
	ID int `json:"id"`
	SiteID int `json:"site_id"`
	StatusCode int `json:"status_code"`
	ResponseTimeMS int64 `json:"response_time_ms"`
	IsAvailable bool `json:"is_available"`
	CheckedAt time.Time `json:"checked_at"`

}