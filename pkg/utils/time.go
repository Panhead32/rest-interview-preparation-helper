package utils

import "time"

func NowTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
