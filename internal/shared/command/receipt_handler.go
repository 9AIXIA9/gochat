package command

import "context"

var _ ReceiptHandler = ReceiptHandlerFunc(nil)

type ReceiptHandler interface {
	Handle(ctx context.Context, r Receipt) error
}

type ReceiptHandlerFunc func(ctx context.Context, r Receipt) error

func (h ReceiptHandlerFunc) Handle(ctx context.Context, r Receipt) error {
	return h(ctx, r)
}
