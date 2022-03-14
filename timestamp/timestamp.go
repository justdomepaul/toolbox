package timestamp

import (
	"math"
	"strconv"
	"time"
)

var timeNow = time.Now

// must use to second timestamp

// GetTimestamp method
func GetTimestamp(timestamp string) time.Time {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(i, 0).UTC()
}

// UTC +8 location

// GetUTC8Time method
func GetUTC8Time(timestamp int64) time.Time {
	return time.Unix(timestamp/int64(time.Microsecond), 0).In(time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds())))
}

// GetUTC8NowTime method
func GetUTC8NowTime() time.Time {
	return time.Unix(timeNow().Unix(), 0).In(time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds())))
}

// GetBeginningOfDay method
func GetBeginningOfDay() int64 {
	currentNow := timeNow()
	zone := time.FixedZone("Asia/Taipei", int((0 * time.Hour).Seconds()))
	utcTime := time.Date(currentNow.Year(), currentNow.Month(), currentNow.Day(), 0, 0, 0, 0, zone)
	return utcTime.UnixNano() / int64(time.Millisecond)
}

// GetUTC8BeginningOfDay method
func GetUTC8BeginningOfDay() int64 {
	currentNow := timeNow()
	zone := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))
	utcTime := time.Date(currentNow.Year(), currentNow.Month(), currentNow.Day(), 0, 0, 0, 0, zone)
	return utcTime.UnixNano() / int64(time.Millisecond)
}

// must use to second timestamp

// GetUnixTimestamp method
func GetUnixTimestamp(timestamp float64) time.Time {
	sec, dec := math.Modf(timestamp)
	return time.Unix(int64(sec), int64(dec*(1e9))).UTC()
}

// blobAsBigint(timestampAsBlob(event_time)) cassandra get millisecond timestamp

// GetNowTimestamp method
func GetNowTimestamp() int64 {
	return timeNow().UnixNano() / int64(time.Millisecond)
}

// GetNowSecondTimestamp method
func GetNowSecondTimestamp() int64 {
	return timeNow().UnixNano() / int64(time.Second)
}

// GetTimeFromTodayBegin method
func GetTimeFromTodayBegin() int64 {
	return GetNowTimestamp() - GetBeginningOfDay()
}

// GetTimeBucket method
func GetTimeBucket(timestamp int64) int64 {
	const mask uint64 = 0xFFFFFFFFF8000000
	return int64(uint64(timestamp) & mask)
}
