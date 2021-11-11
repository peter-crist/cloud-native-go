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
			m.Unlock()
		}()

		log.Println(threshold.Format(time.RFC3339))
		log.Println(time.Now().Format(time.RFC3339))
		if time.Now().Before(threshold) {
			log.Printf("❌ You are being rate limited. Returning most recent results.")
			return result, err
		}

		result, err = circuit(ctx)
		log.Printf("✅ New result received.")
		threshold = time.Now().Add(d)
		return result, err
	}
}
