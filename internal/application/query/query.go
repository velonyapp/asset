package query

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

type Query[R any] interface {
	resultType() R
}

type Handler[Q, R any] interface {
	Handle(context.Context, Q) (R, error)
}

type erasedHandler func(context.Context, any) (any, error)

var handlers sync.Map

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func Register[R any, Q Query[R]](handler Handler[Q, R]) {
	t := typeOf[Q]()

	handlers.Store(t, erasedHandler(func(ctx context.Context, cmd any) (any, error) {
		return handler.Handle(ctx, cmd.(Q))
	}))
}

func Send[R any, Q Query[R]](ctx context.Context, cmd Q) (R, error) {
	var zero R

	raw, ok := handlers.Load(typeOf[Q]())
	if !ok {
		return zero, fmt.Errorf("no handler registered for %T", cmd)
	}

	result, err := raw.(erasedHandler)(ctx, cmd)
	if err != nil {
		return zero, err
	}

	return result.(R), nil
}
