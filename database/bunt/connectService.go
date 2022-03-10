package bunt

import (
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"io"
)

// ISession interface
type ISession interface {
	Close() error
	Save(wr io.Writer) error
	Load(rd io.Reader) error
	CreateIndex(name, pattern string, less ...func(a, b string) bool) error
	ReplaceIndex(name, pattern string, less ...func(a, b string) bool) error
	CreateSpatialIndex(name, pattern string, rect func(item string) (min, max []float64)) error
	ReplaceSpatialIndex(name, pattern string, rect func(item string) (min, max []float64)) error
	DropIndex(name string) error
	Indexes() ([]string, error)
	ReadConfig(config *buntdb.Config) error
	SetConfig(config buntdb.Config) error
	Shrink() error
	View(fn func(tx *buntdb.Tx) error) error
	Update(fn func(tx *buntdb.Tx) error) error
	Begin(writable bool) (*buntdb.Tx, error)
}

// NewSession method
var NewSession = func(path string) (ISession, error) {
	return buntdb.Open(path)
}

func NewExtendBuntDatabase(logger *zap.Logger) (ISession, func(), error) {
	session, err := NewSession(":memory:")
	if err != nil {
		return nil, nil, err
	}
	logger.Info("Bunt init complete", zap.String("system", "Database"))

	return session,
		func() {
			defer errorhandler.PanicErrorHandler("Core", "")
			if err := session.Close(); err != nil {
				panic(errorhandler.NewErrDBDisconnection(err))
			}
		}, nil
}
