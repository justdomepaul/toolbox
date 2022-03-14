package shorten

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ShortenSuite struct {
	suite.Suite
	url string
}

func (suite *ShortenSuite) SetupSuite() {
	suite.url = "https://www.amd.com/zh-hant/where-to-buy/radeon-rx-6000-series-graphics"
}

func (suite *ShortenSuite) TestOutput() {
	input := "https://flashaim.tv/"
	outputDomain := "https://matrix.tw/"
	uid, err := uuid.NewRandom()
	suite.NoError(err)
	result := String(input, uid.String())
	fmt.Printf("Input URL: %s, \nOutput short key: %s \n", input, outputDomain+result)
}

func (suite *ShortenSuite) TestOutputString() {
	input := "https://flashaim.tv/"
	outputDomain := "https://matrix.tw/"

	layout := "2006-01-02T15:04:05.000Z"
	str := "2021-08-25T11:45:26.371Z"
	ti, err := time.Parse(layout, str)
	suite.NoError(err)
	defer gostub.StubFunc(&now, ti).Reset()

	result := Timestamp()
	fmt.Printf("Input URL: %s, \nOutput short key: %s \n", input, outputDomain+result)
}

func (suite *ShortenSuite) TestString() {
	result := String(suite.url, "255dd155-2050-4e58-ac39-62ec524b8ce2")
	suite.Equal("Ak8gvVEpLxJaVi1ntf7fUpRwxL", result)
}

func (suite *ShortenSuite) TestTimestamp() {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2021-08-25T11:45:26.371Z"
	ti, err := time.Parse(layout, str)
	suite.NoError(err)
	defer gostub.StubFunc(&now, ti).Reset()
	result := Timestamp()
	suite.Equal("anQkM7JHV9", result)
}

func (suite *ShortenSuite) TestMaxTimestamp() {
	suite.Equal("2032-08-20T19:43:03.67Z", MaxTimestamp(35, 10*time.Millisecond).UTC().Format(time.RFC3339Nano))
	suite.Equal("2056-08-03T19:53:47.775Z", MaxTimestamp(40, time.Millisecond).UTC().Format(time.RFC3339Nano))
	suite.Equal("2035-09-08T07:57:31.1103Z", MaxTimestamp(42, 100*time.Microsecond).UTC().Format(time.RFC3339Nano))
	suite.Equal("2110-12-12T02:56:07.10655Z", MaxTimestamp(48, 10*time.Microsecond).UTC().Format(time.RFC3339Nano))
	suite.Equal("FSb3FGK", timestamp(40, 1))
	suite.Equal("zmM9z4F", timestamp(42, 1))
	suite.Equal("2zXhJy7U", timestamp(42, 1<<42-1))
	suite.Equal("zmM9z4E", timestamp(42, 1<<42-2))
}

func (suite *ShortenSuite) TestTimestamp42() {
	// basic
	t0, err := time.Parse(time.RFC3339Nano, "2021-10-13T02:49:57.000000000Z")
	suite.NoError(err)
	defer gostub.StubFunc(&now, t0).Reset()
	s0 := Timestamp42()
	suite.Equal("28fZqLs", s0)

	// under resolution
	t1 := t0.Add(99 * time.Microsecond)
	defer gostub.StubFunc(&now, t1).Reset()
	s1 := Timestamp42()
	suite.Equal(s0, s1)

	// exceed
	t2 := t0.Add(100 * time.Microsecond)
	defer gostub.StubFunc(&now, t2).Reset()
	s2 := Timestamp42()
	suite.NotEqual(s0, s2)
	suite.Equal("GaFc5cB", s2)

	// under resolution, again
	t3 := t0.Add(101 * time.Microsecond)
	defer gostub.StubFunc(&now, t3).Reset()
	s3 := Timestamp42()
	suite.Equal(s2, s3)
}

func TestShortenSuite(t *testing.T) {
	suite.Run(t, new(ShortenSuite))
}

func BenchmarkShortenString(b *testing.B) {
	url := "https://www.amd.com/zh-hant/where-to-buy/radeon-rx-6000-series-graphics"
	for i := 0; i < b.N; i++ {
		String(url, "e0dba740-fc4b-4977-872c-d360239e6b1a")
	}
}

func BenchmarkShortenTimestamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Timestamp()
	}
}
