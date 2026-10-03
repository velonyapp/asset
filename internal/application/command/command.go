package command

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

type Command[R any] interface {
	resultType() R
}

type Handler[C, R any] interface {
	Handle(context.Context, C) (R, error)
}

type erasedHandler func(context.Context, any) (any, error)

var handlers sync.Map

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func Register[R any, C Command[R]](handler Handler[C, R]) {
	t := typeOf[C]()

	handlers.Store(t, erasedHandler(func(ctx context.Context, cmd any) (any, error) {
		return handler.Handle(ctx, cmd.(C))
	}))
}

func Send[R any, C Command[R]](ctx context.Context, cmd C) (R, error) {
	var zero R

	raw, ok := handlers.Load(typeOf[C]())
	if !ok {
		return zero, fmt.Errorf("no handler registered for %T", cmd)
	}

	result, err := raw.(erasedHandler)(ctx, cmd)
	if err != nil {
		return zero, err
	}

	return result.(R), nil
}
