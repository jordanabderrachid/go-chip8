package timer

import "time"

const (
	Frenquency int64 = 60 // Hz
)

type Timer struct {
}

func (Timer) Run(b *byte) {
	ticker := time.NewTicker(time.Duration(int64(time.Second) / Frenquency))
	for {
		select {
		case <-ticker.C:
			if *b > 0 {
				*b--
			}
		}
	}
}
