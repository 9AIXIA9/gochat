package websocket

import "context"

type Validator interface {
	Validate(ctx context.Context, model any) (string, error)
}
