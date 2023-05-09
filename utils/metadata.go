package utils

import (
	"context"
	"github.com/justdomepaul/toolbox/generic"
)

func SetID[T, E generic.ByteSeq](ctx context.Context, key T, id E) context.Context {
	return context.WithValue(ctx, key, id)
}

func GetID[T generic.ByteSeq, E generic.ByteSeq](ctx context.Context, key T) E {
	return ctx.Value(key).(E)
}

func SetClaim[T generic.ByteSeq](ctx context.Context, key T, claim interface{}) context.Context {
	return context.WithValue(ctx, key, claim)
}

func GetClaim[T generic.ByteSeq](ctx context.Context, key T) interface{} {
	return ctx.Value(key)
}
