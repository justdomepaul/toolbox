package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/errorhandler"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	newClient = pubsub.NewClient
)

// ISession interface
type ISession interface {
	Close() error
	Snapshot(id string) *pubsub.Snapshot
	Snapshots(ctx context.Context) *pubsub.SnapshotConfigIterator
	Subscription(id string) *pubsub.Subscription
	SubscriptionInProject(id, projectID string) *pubsub.Subscription
	Subscriptions(ctx context.Context) *pubsub.SubscriptionIterator
	CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error)
	CreateTopic(ctx context.Context, topicID string) (*pubsub.Topic, error)
	CreateTopicWithConfig(ctx context.Context, topicID string, tc *pubsub.TopicConfig) (*pubsub.Topic, error)
	Topic(id string) *pubsub.Topic
	TopicInProject(id, projectID string) *pubsub.Topic
	Topics(ctx context.Context) *pubsub.TopicIterator
}

// NewSession method
var NewSession = func(ctx context.Context, opt config.PubSub) (ISession, error) {
	options := make([]option.ClientOption, 0)
	if opt.EndPoint != "" {
		options = append(options, option.WithEndpoint(opt.EndPoint))
	}
	if opt.WithoutAuthentication {
		options = append(options, option.WithoutAuthentication())
	}
	if opt.GRPCInsecure {
		options = append(options, option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	}
	return newClient(ctx, fmt.Sprintf(`projects/%s`, opt.ProjectID), options...)
}

func NewExtendPubSubDatabase(logger *zap.Logger, opt config.PubSub) (ISession, func(), error) {
	session, err := NewSession(context.Background(), opt)
	if err != nil {
		return nil, nil, err
	}
	logger.Info("PubSub init complete", zap.String("system", "Database"))

	return session,
		func() {
			defer errorhandler.PanicErrorHandler("Core", "")
			session.Close()
		}, nil
}
