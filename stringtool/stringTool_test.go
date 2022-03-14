package stringtool

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type StringToolSuite struct {
	suite.Suite
}

func (suite *StringToolSuite) TestStringJoin() {
	suite.Equal("Hello world", StringJoin("Hello", " ", "world"))
}

func (suite *StringToolSuite) TestConvertFirstUpper() {
	suite.Equal("Helloworld", ConvertFirstUpper("helloworld"))
}

func (suite *StringToolSuite) TestConvertFirstLower() {
	suite.Equal("helloworld", ConvertFirstLower("Helloworld"))
}

func TestStringToolSuite(t *testing.T) {
	suite.Run(t, new(StringToolSuite))
}

func Benchmark_StringJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringJoin("Hello", " ", "world")
	}
}

func Benchmark_ConvertFirstUpper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConvertFirstUpper("helloworld")
	}
}

func Benchmark_ConvertFirstLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConvertFirstLower("Helloworld")
	}
}
