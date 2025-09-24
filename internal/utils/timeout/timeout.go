package timeout

import (
	"context"
	"errors"
	"gochat/internal/types"
)

type (
	CtxFnWithResponse[response any] func(context.Context) (response, error)

	NonCtxFnWithResponse[response any] func() (response, error)

	// CtxFn 是一个检查上下文的函数类型
	CtxFn func(context.Context) error

	// NonCtxFn 是一个无上下文检查的函数类型
	NonCtxFn func() error
)

func ConvertAndExecuteWithResponse[response any](ctx context.Context, fn NonCtxFnWithResponse[response]) (response, error) {
	fnn := ConvertWithResponse(ctx, fn)
	return fnn(ctx)
}

func ConvertWithResponse[response any](ctx context.Context, fn NonCtxFnWithResponse[response]) CtxFnWithResponse[response] {
	return func(context.Context) (response, error) {
		var resp response
		fnn := Convert(ctx, func() error {
			var err error
			resp, err = fn()
			return err
		})

		err := fnn(ctx)
		return resp, err
	}
}

func ConvertAndExecute(ctx context.Context, fn NonCtxFn) error {
	fnn := Convert(ctx, fn)
	return fnn(ctx)
}

func Convert(ctx context.Context, fn NonCtxFn) CtxFn {
	return func(context.Context) error {
		err := Check(ctx)
		if err != nil {
			return err
		}
		return fn()
	}
}

// Check 检查上下文是否已取消或超时
func Check(ctx context.Context) error {
	if IsCanceled(ctx.Err()) {
		return types.ErrCanceled
	}
	if IsTimeout(ctx.Err()) {
		return types.ErrTimeout
	}

	return ctx.Err()
}

// IsCanceledOrTimeout 错误是上下文被取消或者是超时
func IsCanceledOrTimeout(err error) bool {
	return IsCanceled(err) || IsTimeout(err)
}

// IsCanceled 判断错误是否由上下文取消引起
func IsCanceled(err error) bool {
	return errors.Is(err, types.ErrCanceled) || errors.Is(err, context.Canceled)
}

// IsTimeout 判断错误是否由上下文超时引起
func IsTimeout(err error) bool {
	return errors.Is(err, types.ErrTimeout) || errors.Is(err, context.DeadlineExceeded)
}
