package redis

import (
	"context"
	"fmt"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

// ISession interface
type ISession interface {
	redis.Cmdable
	WithTimeout(timeout time.Duration) *redis.Client
	Conn() *redis.Conn
	Do(ctx context.Context, args ...interface{}) *redis.Cmd
	Process(ctx context.Context, cmd redis.Cmder) error
	Options() *redis.Options
	PoolStats() *redis.PoolStats
	Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error)
	Pipeline() redis.Pipeliner
	TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error)
	TxPipeline() redis.Pipeliner
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
	PSubscribe(ctx context.Context, channels ...string) *redis.PubSub
	SSubscribe(ctx context.Context, channels ...string) *redis.PubSub
	String() string
	Close() error
	AddHook(hook redis.Hook)
}

// NewSession method
var NewSession = func(opt config.Redis) (ISession, error) {
	rdsOpt := &redis.Options{
		Addr: fmt.Sprintf("%s:%s", opt.RedisHost, opt.RedisPort),
		DB:   0, // use default DB
	}
	if opt.RedisClientName != "" {
		rdsOpt.ClientName = opt.RedisClientName
	}
	if opt.RedisUsername != "" {
		rdsOpt.Username = opt.RedisUsername
	}
	if opt.RedisPassword != "" {
		rdsOpt.Password = opt.RedisPassword
	}
	if opt.RedisPoolSize > 0 {
		rdsOpt.PoolSize = opt.RedisPoolSize
	}
	return redis.NewClient(rdsOpt), nil
}

func NewExtendRedisDatabase(logger *zap.Logger, opt config.Redis) (ISession, func(), error) {
	session, err := NewSession(opt)
	if err != nil {
		return nil, nil, err
	}
	logger.Info("Redis init complete", zap.String("system", "Database"))

	return session,
		func() {
			defer errorhandler.PanicErrorHandler("Core", "")
			if err := session.Close(); err != nil {
				panic(errorhandler.NewErrDBDisconnection(err))
			}
		}, nil
}
