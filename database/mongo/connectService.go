package mongo

import (
	"context"
	"fmt"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/errorhandler"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"time"
)

var (
	newClient = mongo.Connect
)

// ISession interface
type ISession interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	StartSession(opts ...*options.SessionOptions) (mongo.Session, error)
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
	ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (mongo.ListDatabasesResult, error)
	ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error)
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
	UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error
	Watch(ctx context.Context, pipeline interface{},
		opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	NumberSessionsInProgress() int
	Timeout() *time.Duration
}

// NewSession method
var NewSession = func(ctx context.Context, opt config.Mongo) (ISession, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	authentication := ""
	database := ""
	permission := "?retryWrites=true&w=majority"
	if opt.MongoUsername != "" && opt.MongoPassword != "" {
		authentication = fmt.Sprintf("%s:%s@", opt.MongoUsername, opt.MongoPassword)
	}
	if opt.MongoDatabase != "" {
		database = fmt.Sprintf("/%s", opt.MongoDatabase)
	}
	if opt.MongoAuthSource {
		permission = "?authSource=admin"
	}

	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("%s://%s%s%s%s", opt.MongoProtocol, authentication, opt.MongoHost, database, permission)).
		SetServerAPIOptions(serverAPIOptions)
	return newClient(ctx, clientOptions)
}

func NewExtendMongoDatabase(logger *zap.Logger, opt config.Mongo) (ISession, func(), error) {
	ctx := context.Background()
	session, err := NewSession(ctx, opt)
	if err != nil {
		return nil, nil, err
	}
	logger.Info("MongoDB init complete", zap.String("system", "Database"))

	return session,
		func() {
			defer errorhandler.PanicErrorHandler("Core", "")
			if err := session.Disconnect(ctx); err != nil {
				panic(errorhandler.NewErrDBDisconnection(err))
			}
		}, nil
}
