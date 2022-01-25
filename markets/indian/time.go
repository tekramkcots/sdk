package indian

import "time"

func GetTimeZone() *time.Location {
	//+5:30 in seconds = 19800
	return time.FixedZone("IST", 19800)
}
