package array

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type FindSuite struct {
	suite.Suite
}

func (suite *FindSuite) TestFindString() {
	input := []string{"testA", "testB", "testC"}
	index, exist := Find(input, "testB")
	suite.True(exist)
	suite.Equal(1, index)
}

func (suite *FindSuite) TestFindStringNotExit() {
	input := []string{"testA", "testB", "testC"}
	index, exist := Find(input, "testD")
	suite.False(exist)
	suite.Equal(-1, index)
}

func (suite *FindSuite) TestFindInt() {
	input := []int{0, 5, 11, 13, 255}
	index, exist := Find(input, 13)
	suite.True(exist)
	suite.Equal(3, index)
}

func (suite *FindSuite) TestFindIntNotExist() {
	input := []int{0, 5, 11, 13, 255}
	index, exist := Find(input, 205)
	suite.False(exist)
	suite.Equal(-1, index)
}

func TestFindSuite(t *testing.T) {
	suite.Run(t, new(FindSuite))
}
