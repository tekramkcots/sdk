package indian

import "time"

func GetTimeZone() *time.Location {
	//+5:30 in seconds = 19800
	return time.FixedZone("Asia/Kolkata", 19800)
}
