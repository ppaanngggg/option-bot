package utils

import "time"

var TZNewYork *time.Location

func init() {
	TZNewYork, _ = time.LoadLocation("America/New_York")
}
