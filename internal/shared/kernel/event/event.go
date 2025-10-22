package event

import (
	"gochat/internal/shared/kernel"
)

type Event interface {
	kernel.Serializer
	StandardEvent
}
