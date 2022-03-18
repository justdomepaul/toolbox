package shutdown

import (
	"github.com/stretchr/testify/suite"
	"os"
	"reflect"
	"syscall"
	"testing"
	"time"
)

type shutdown struct {
	suite.Suite
}

func (suite *shutdown) TestWithQuitOption() {
	suite.Equal("chan os.Signal", reflect.TypeOf(NewShutdown(WithQuit(make(chan os.Signal))).quit).String())
}

func (suite *shutdown) TestWithDoneOption() {
	suite.Equal("chan bool", reflect.TypeOf(NewShutdown(WithDone(make(chan bool))).done).String())
}

func (suite *shutdown) TestWithServerTimeoutOption() {
	suite.Equal("time.Duration", reflect.TypeOf(NewShutdown(WithServerTimeout(5*time.Second)).serverTimeout).String())
}

func (suite *shutdown) TestWithEndTaskOption() {
	suite.Equal("func()", reflect.TypeOf(NewShutdown(WithEndTask(func() {})).endTask).String())
}

func TestShutdown(t *testing.T) {
	suite.Run(t, new(shutdown))
}

func Test_Shutdown(t *testing.T) {
	t.Run("test Shutdown", func(t *testing.T) {
		quit := make(chan os.Signal)
		defer close(quit)
		done := make(chan bool)
		defer close(done)

		go func() {
			(&Shutdown{
				quit: quit,
				done: done,
				endTask: func() {
					t.Log("endTask")
				},
			}).Shutdown()
		}()
		quit <- syscall.SIGTERM
		<-done
	})
}
