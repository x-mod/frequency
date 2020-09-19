package frequency

import (
	"time"

	"github.com/juju/ratelimit"
)

type Frequency struct {
	blocked bool
	limits  []*ratelimit.Bucket
}

type Option func(*Frequency)

func Limit(duration time.Duration, max int64, cap int64) Option {
	return func(f *Frequency) {
		if max > 0 {
			bucket := ratelimit.NewBucketWithQuantum(duration, cap, max)
			f.limits = append(f.limits, bucket)
		} else {
			f.blocked = true
		}
	}
}

func Second(max int64) Option {
	return Limit(time.Second, max, max)
}

func Minute(max int64) Option {
	return Limit(time.Minute, max, max)
}

func Hour(max int64) Option {
	return Limit(time.Hour, max, max)
}

func Day(max int64) Option {
	return Limit(time.Hour*24, max, max)
}

func New(opts ...Option) *Frequency {
	freq := &Frequency{}
	for _, opt := range opts {
		opt(freq)
	}
	return freq
}

func (freq *Frequency) IsBlocked() bool {
	return freq.blocked
}

func (freq *Frequency) ReserveN(n int64) (time.Duration, bool) {
	if freq.blocked {
		return 0, false
	}
	var r time.Duration
	for _, limit := range freq.limits {
		d := limit.Take(n)
		if d > r {
			r = d
		}
	}
	return r, true
}

func (freq *Frequency) WaitN(n int64) {
	for _, limit := range freq.limits {
		limit.Wait(n)
	}
}
