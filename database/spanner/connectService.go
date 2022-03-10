package spanner

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/errorhandler"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"time"
)

var (
	newClient = spanner.NewClient
)

// ISession interface
type ISession interface {
	Close()
	Single() *spanner.ReadOnlyTransaction
	ReadOnlyTransaction() *spanner.ReadOnlyTransaction
	BatchReadOnlyTransaction(ctx context.Context, tb spanner.TimestampBound) (*spanner.BatchReadOnlyTransaction, error)
	BatchReadOnlyTransactionFromID(tid spanner.BatchReadOnlyTransactionID) *spanner.BatchReadOnlyTransaction
	ReadWriteTransaction(ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error) (commitTimestamp time.Time, err error)
	Apply(ctx context.Context, ms []*spanner.Mutation, opts ...spanner.ApplyOption) (commitTimestamp time.Time, err error)
	PartitionedUpdate(ctx context.Context, statement spanner.Statement) (count int64, err error)
	PartitionedUpdateWithOptions(ctx context.Context, statement spanner.Statement, opts spanner.QueryOptions) (count int64, err error)
}

// NewSession method
var NewSession = func(ctx context.Context, opt config.Spanner) (ISession, error) {
	options := make([]option.ClientOption, 0)
	if opt.EndPoint != "" {
		options = append(options, option.WithEndpoint(opt.EndPoint))
	}
	if opt.WithoutAuthentication {
		options = append(options, option.WithoutAuthentication())
	}
	if opt.GRPCInsecure {
		options = append(options, option.WithGRPCDialOption(grpc.WithInsecure()))
	}
	return newClient(ctx, fmt.Sprintf(`projects/%s/instances/%s/databases/%s`, opt.ProjectID, opt.Instance, opt.Database), options...)
}

func NewExtendSpannerDatabase(logger *zap.Logger, opt config.Spanner) (ISession, func(), error) {
	session, err := NewSession(context.Background(), opt)
	if err != nil {
		return nil, nil, err
	}
	logger.Info("Spanner init complete", zap.String("system", "Database"))

	return session,
		func() {
			defer errorhandler.PanicErrorHandler("Core", "")
			session.Close()
		}, nil
}
