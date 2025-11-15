package kernel

import (
	"context"
)

// UseCase 一个用例对应一个请求(Input)与一个响应(Output)。
type UseCase[
	Input Validatable,
	Output any,
] interface {
	Execute(ctx context.Context, input Input) (Output, error)
}

// NoInput 无入参时使用。
type NoInput struct{}

func (*NoInput) Validate() error { return nil }

// NoOutput 无返回体时使用。
type NoOutput struct{}
