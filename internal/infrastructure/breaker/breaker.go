package breaker

import (
	myErrors "gochat/internal/shared/errors"

	"github.com/sony/gobreaker/v2"
	"go.uber.org/zap"
)

func NewBreaker[T any](conf *Config) *gobreaker.CircuitBreaker[T] {
	return gobreaker.NewCircuitBreaker[T](gobreaker.Settings{
		Name:         conf.Name,
		MaxRequests:  conf.MaxRequests,
		Interval:     conf.Interval,
		BucketPeriod: conf.BucketPeriod,
		Timeout:      conf.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if counts.Requests < conf.ErrorThreshold {
				return false
			}
			var errorRate float64
			if counts.Requests > 0 {
				errorRate = float64(counts.TotalFailures) / float64(counts.Requests)
			}
			return errorRate >= conf.ErrorRate
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			zap.L().Info(
				"circuit breaker state changed",
				zap.String("name", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()),
			)
		},
		IsSuccessful: func(err error) bool {
			return err == nil || myErrors.IsBusinessError(err)
		},
	})
}
