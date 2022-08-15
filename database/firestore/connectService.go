package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/errorhandler"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

var (
	newClient = firestore.NewClient
)

// ISession interface
type ISession interface {
	Close() error
	Collection(path string) *firestore.CollectionRef
	Doc(path string) *firestore.DocumentRef
	CollectionGroup(collectionID string) *firestore.CollectionGroupRef
	GetAll(ctx context.Context, docRefs []*firestore.DocumentRef) (_ []*firestore.DocumentSnapshot, err error)
	Collections(ctx context.Context) *firestore.CollectionIterator
	Batch() *firestore.WriteBatch
	RunTransaction(ctx context.Context, f func(context.Context, *firestore.Transaction) error, opts ...firestore.TransactionOption) (err error)
}

// NewSession method
var NewSession = func(ctx context.Context, opt config.Firestore) (ISession, error) {
	options := make([]option.ClientOption, 0)
	if opt.EndPoint != "" {
		options = append(options, option.WithEndpoint(opt.EndPoint))
	}
	return newClient(ctx, opt.ProjectID, options...)
}

func NewExtendFirestoreDatabase(logger *zap.Logger, opt config.Firestore) (ISession, func(), error) {
	session, err := NewSession(context.Background(), opt)
	if err != nil {
		return nil, nil, err
	}
	logger.Info("Firestore init complete", zap.String("system", "Database"))

	return session,
		func() {
			defer errorhandler.PanicErrorHandler("Core", "")
			if err := session.Close(); err != nil {
				panic(errorhandler.NewErrDBDisconnection(err))
			}
		}, nil
}
