package circuitbreaker

import (
	"errors"
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var (
		consecutiveFailures int = 0
		lastAttempt             = time.Now()
		m                   sync.RWMutex
	)

	return func(ctx context.Context) (string, error) {
		m.RLock()
		excessiveAttempts := consecutiveFailures - int(failureThreshold)
		log.Printf("%d consecutive failure(s). %d attempts remaining", consecutiveFailures, excessiveAttempts)
		if excessiveAttempts >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << excessiveAttempts)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				log.Println("💥 Breaker tripped! Service unreachable, please try again later.")
				return "", errors.New("service unreachable")
			} else {
				log.Println("✅ Breaker reset - ready for more connection attempts.")
			}
		}

		m.RUnlock()

		response, err := circuit(ctx)

		m.Lock()
		defer m.Unlock()

		lastAttempt = time.Now()
		if err != nil {
			consecutiveFailures++
			return response, err
		}

		consecutiveFailures = 0
		return response, nil
	}

}
