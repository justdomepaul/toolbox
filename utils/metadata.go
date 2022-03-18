package utils

import (
	"context"
	"github.com/justdomepaul/toolbox/generic"
)

func SetClientID[T, E generic.ByteSeq](ctx context.Context, key T, id E) context.Context {
	return context.WithValue(ctx, key, id)
}

func GetClientID[T generic.ByteSeq, E generic.ByteSeq](ctx context.Context, key T) E {
	return ctx.Value(key).(E)
}
