package loggingadmin

import (
	"cloud.google.com/go/logging/logadmin"
	"context"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/errorhandler"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

var (
	newClient = logadmin.NewClient
)

// ISession interface
type ISession interface {
	Close() error
	DeleteLog(ctx context.Context, logID string) error
	Entries(ctx context.Context, opts ...logadmin.EntriesOption) *logadmin.EntryIterator
	Logs(ctx context.Context) *logadmin.LogIterator
	CreateMetric(ctx context.Context, m *logadmin.Metric) error
	DeleteMetric(ctx context.Context, metricID string) error
	Metric(ctx context.Context, metricID string) (*logadmin.Metric, error)
	UpdateMetric(ctx context.Context, m *logadmin.Metric) error
	Metrics(ctx context.Context) *logadmin.MetricIterator
	ResourceDescriptors(ctx context.Context) *logadmin.ResourceDescriptorIterator
	CreateSink(ctx context.Context, sink *logadmin.Sink) (*logadmin.Sink, error)
	CreateSinkOpt(ctx context.Context, sink *logadmin.Sink, opts logadmin.SinkOptions) (*logadmin.Sink, error)
	DeleteSink(ctx context.Context, sinkID string) error
	Sink(ctx context.Context, sinkID string) (*logadmin.Sink, error)
	UpdateSink(ctx context.Context, sink *logadmin.Sink) (*logadmin.Sink, error)
	UpdateSinkOpt(ctx context.Context, sink *logadmin.Sink, opts logadmin.SinkOptions) (*logadmin.Sink, error)
	Sinks(ctx context.Context) *logadmin.SinkIterator
}

// NewSession method
var NewSession = func(ctx context.Context, opt config.Spanner, opts ...option.ClientOption) (ISession, error) {
	return newClient(ctx, opt.ProjectID, opts...)
}

func NewExtendLoggingAdmin(logger *zap.Logger, opt config.Spanner) (ISession, func(), error) {
	session, err := NewSession(context.Background(), opt)
	if err != nil {
		return nil, nil, err
	}
	logger.Info("logging admin init complete", zap.String("system", "Database"))

	return session,
		func() {
			defer errorhandler.PanicErrorHandler("Core", "")
			if err := session.Close(); err != nil {
				panic(err)
			}
		}, nil
}
