package main

import (
	"userms/internal/router"
)

func tokenPeriodTimestamp(secs int64) router.OptsFunc {
	return func(opts *router.Options) {
		opts.TokenPeriodTimestamp = secs
	}
}
