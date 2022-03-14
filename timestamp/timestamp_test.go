package timestamp

import (
	"fmt"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TimestampSuite struct {
	suite.Suite
}

func (suite *TimestampSuite) TestGetTimestamp() {
	suite.Equal("2019-01-03 09:39:57 +0000 UTC", GetTimestamp("1546508397").String())
}

func (suite *TimestampSuite) TestGetTimestampFormatError() {
	suite.Panics(func() {
		GetTimestamp("15465asdf08397")
	})
}

func (suite *TimestampSuite) TestGetUTC8Time() {
	suite.Equal("2018-09-19 15:45:57 +0800 Asia/Taipei", GetUTC8Time(1537343157984).String())
}

func (suite *TimestampSuite) TestGetUnixTimestamp() {
	suite.Equal("2019-01-03 09:39:57 +0000 UTC", GetUnixTimestamp(1546508397).String())
}

func (suite *TimestampSuite) TestGetUTC8NowTime() {
	defer gostub.Stub(&timeNow, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()

	suite.Equal(int64(1550836520000), GetUTC8NowTime().UnixNano()/int64(time.Millisecond))
	suite.Equal("2019/02/22 19:55:20", GetUTC8NowTime().Format("2006/01/02 15:04:05"))
}

func (suite *TimestampSuite) TestGetBeginningOfDay() {
	defer gostub.Stub(&timeNow, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	suite.Equal(int64(1550793600000), GetBeginningOfDay())
}

func (suite *TimestampSuite) TestGetUTC8BeginningOfDay() {
	defer gostub.Stub(&timeNow, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	suite.Equal(int64(1550764800000), GetUTC8BeginningOfDay())
}

func (suite *TimestampSuite) TestGetNowTimestamp() {
	defer gostub.Stub(&timeNow, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	suite.Equal(int64(1550836520000), GetNowTimestamp())
}

func (suite *TimestampSuite) TestGetNowSecondTimestamp() {
	defer gostub.Stub(&timeNow, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	suite.Equal(int64(1550836520), GetNowSecondTimestamp())
}

func (suite *TimestampSuite) TestGetTimeFromTodayBegin() {
	defer gostub.Stub(&timeNow, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	suite.Equal(int64(42920000), GetTimeFromTodayBegin())
}

func (suite *TimestampSuite) TestGetTimeBucket() {
	defer gostub.Stub(&timeNow, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	suite.Equal(int64(1537329856512), GetTimeBucket(1537343157984))
}

func TestTimestampSuite(t *testing.T) {
	suite.Run(t, new(TimestampSuite))
}

func Benchmark_GetTimestamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetTimestamp("1546508397")
	}
}

func Benchmark_GetUTC8Time(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetUTC8Time(1537343157984)
	}
}

func Benchmark_GetUnixTimestamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetUnixTimestamp(1546508397)
	}
}

func Benchmark_GetTimeBucket(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetTimeBucket(1537343157984)
	}
}

func ExampleGetTimeBucket() {
	result := GetTimeBucket(1537343157984)
	fmt.Println(result)
	// Output: 1537329856512
}
