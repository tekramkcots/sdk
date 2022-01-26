package indian

import "time"

func GetTimeZone() *time.Location {
	//+5:30 in seconds = 19800
	return time.FixedZone("IST", 19800)
}

func MarketStartTime() time.Time {
	return time.Date(2021, time.December, 26, 9, 15, 0, 0, GetTimeZone())
}
