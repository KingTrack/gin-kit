package context

import (
	"context"
	"fmt"
)

type Key struct{}

type Value string

func (v Value) ToString() string {
	return string(v)
}

func (v Value) IsEmpty() bool {
	return len(v) == 0
}

const (
	Default Value = ""
)

var (
	traceIDKey   = &Key{}
	namespaceKey = &Key{}
)

func NewBackground(v Value) context.Context {
	return context.WithValue(context.Background(), namespaceKey, v.ToString())
}

func WithTraceID(ctx context.Context, v Value) context.Context {
	return context.WithValue(ctx, traceIDKey, v)
}

func WithNamespace(ctx context.Context, v Value) context.Context {
	return context.WithValue(ctx, namespaceKey, v)
}

func GetNamespace(ctx context.Context) Value {
	v, ok := ctx.Value(namespaceKey).(string)
	if ok {
		return Value(v)
	}
	return Default
}

func GetTraceID(ctx context.Context) Value {
	v, ok := ctx.Value(traceIDKey).(string)
	if ok {
		return Value(v)
	}
	return Default
}

type ResourceName string

func (r ResourceName) ToString() string {
	return string(r)
}

func GetResourceName(ctx context.Context, name string) ResourceName {
	v := GetNamespace(ctx)
	if v.IsEmpty() {
		return ResourceName(name)
	}
	return ResourceName(fmt.Sprintf("%s.%s", v, name))
}
