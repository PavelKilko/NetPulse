package models

import "time"

type Metrics struct {
	URLID        uint      `bson:"url_id"`
	ResponseTime int       `bson:"response_time"`
	StatusCode   int       `bson:"status_code"`
	Timestamp    time.Time `bson:"timestamp"`
}
