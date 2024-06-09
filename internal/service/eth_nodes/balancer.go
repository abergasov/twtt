package ethnodes

import (
	"context"
	"fmt"
	"time"
	"twtt/internal/logger"
	"twtt/internal/utils"
)

const (
	DefaultNormalRetries   = 100
	DefaultFallbackRetries = 5
)

type BalancerConfig struct {
	NormalURLs    []string      `yaml:"rpc_nodes"`
	FallbackURLs  []string      `yaml:"rpc_fallback_nodes"`
	NormalRetries int           `yaml:"normal_retries,omitempty"`
	Timeout       time.Duration `yaml:"timeout,omitempty"`
}

type Balancer[T any] struct {
	normal          *utils.RoundRobinBalancer[T]
	fallback        *utils.RoundRobinBalancer[T]
	normalRetries   int
	fallbackRetries int
	timeout         *time.Duration
	logger          logger.AppLogger
	requestID       uint64
}

type InitBalancersFunc[T any] func() (*utils.RoundRobinBalancer[T], *utils.RoundRobinBalancer[T])

func NewBalancer[T any](initFn InitBalancersFunc[T], normalRetries, fallbackRetries int, timeout *time.Duration, log logger.AppLogger) *Balancer[T] {
	normal, fallback := initFn()
	balancer := &Balancer[T]{
		normal:          normal,
		fallback:        fallback,
		normalRetries:   DefaultNormalRetries,
		fallbackRetries: DefaultFallbackRetries,
		timeout:         timeout,
		logger:          log,
	}

	if normalRetries > 0 {
		balancer.normalRetries = normalRetries
	}
	if fallbackRetries > 0 {
		balancer.fallbackRetries = fallbackRetries
	}

	return balancer
}

type BalancerRepeatableFuncFactory[E any, T any] func(nodes *utils.RoundRobinBalancer[E]) func() (T, error)

func Run[E any, T any](ctx context.Context, fn BalancerRepeatableFuncFactory[E, T], b *Balancer[E], method string, exitErrors ...error) (T, error) {
	repeater := utils.NewFuncRepeater(fn(b.normal)).
		WithCtx(ctx).
		WithFallback(fn(b.fallback)).
		WithMaxTries(b.normalRetries).
		WithFallbackTries(b.fallbackRetries).
		WithExitErrors(exitErrors...).
		WithErrMsg(fmt.Sprintf("failed to get %s request to RPC node", method))

	if b.timeout != nil {
		repeater.WithTriesTimeout(*b.timeout).WithLinearTimeout(true)
	}

	return repeater.Run()
}
