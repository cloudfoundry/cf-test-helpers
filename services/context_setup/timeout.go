package context_setup

import (
	"time"
)

var TimeoutScale float64

func ScaledTimeout(timeout time.Duration) time.Duration {
	return time.Duration(float64(timeout) * TimeoutScale)
}
