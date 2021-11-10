package debounce

import (
	"context"
	"log"
	"sync"
	"time"
)

type Circuit func(context.Context) (string, error)

func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
	var (
		threshold time.Time
		result    string
		err       error
		m         sync.Mutex
	)

	return func(ctx context.Context) (string, error) {
		m.Lock()

		defer func() {
			threshold = time.Now().Add(d)
			m.Unlock()
		}()

		if time.Now().Before(threshold) {
			log.Printf("❌ You are being rate limited. Returning most recent results.")
			return result, err
		}

		result, err = circuit(ctx)
		log.Printf("✅ New result received.")
		return result, err
	}
}
