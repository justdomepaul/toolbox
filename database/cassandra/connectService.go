package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/justdomepaul/toolbox/config"
	"go.uber.org/zap"
	"strconv"
)

// IClusterConfig interface
type IClusterConfig interface {
	CreateSession() (*gocql.Session, error)
}

// ISession interface
type ISession interface {
	SetConsistency(cons gocql.Consistency)
	SetPageSize(n int)
	SetPrefetch(p float64)
	SetTrace(trace gocql.Tracer)
	Query(stmt string, values ...interface{}) *gocql.Query
	Bind(stmt string, b func(q *gocql.QueryInfo) ([]interface{}, error)) *gocql.Query
	Close()
	Closed() bool
	KeyspaceMetadata(string) (*gocql.KeyspaceMetadata, error)
	ExecuteBatch(batch *gocql.Batch) error
	ExecuteBatchCAS(batch *gocql.Batch, dest ...interface{}) (applied bool, iter *gocql.Iter, err error)
	MapExecuteBatchCAS(batch *gocql.Batch, dest map[string]interface{}) (applied bool, iter *gocql.Iter, err error)
	NewBatch(typ gocql.BatchType) *gocql.Batch
}

var NewCluster = func(config *gocql.ClusterConfig) IClusterConfig {
	return config
}

// NewSession method
var NewSession = func(opt config.Cassandra) (ISession, error) {
	cluster := gocql.NewCluster(opt.CassandraHosts...)
	if err := extendOptions(cluster, opt); err != nil {
		return nil, err
	}
	return NewCluster(cluster).CreateSession()
}

var extendOptions = func(cluster *gocql.ClusterConfig, opt config.Cassandra) error {
	if opt.CassandraForcePassword {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: opt.CassandraUsername,
			Password: opt.CassandraPassword,
		}
	}
	cluster.Compressor = &gocql.SnappyCompressor{}
	cluster.Keyspace = opt.CassandraKeySpace
	cluster.ReconnectInterval = opt.ReconnectInterval
	cluster.Timeout = opt.CassandraTimeout
	cluster.ConnectTimeout = opt.CassandraConnectionTimeout
	cluster.NumConns = opt.CassandraNumConnection
	cluster.DisableInitialHostLookup = opt.DisableInitialHostLookup
	port := opt.CassandraPort
	if port != "" {
		var errAToi error
		cluster.Port, errAToi = strconv.Atoi(port)
		if errAToi != nil {
			return errAToi
		}
	}
	return nil
}

func NewExtendCassandraDatabase(logger *zap.Logger, opt config.Cassandra) (ISession, func(), error) {
	session, err := NewSession(opt)
	if err != nil {
		return nil, nil, err
	}
	logger.Info("CassandraDB init complete", zap.String("system", "Database"))

	return session,
		func() {
			session.Close()
		}, nil
}
