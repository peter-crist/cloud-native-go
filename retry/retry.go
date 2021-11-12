package retry

import (
	"context"
	"log"
	"time"
)

type Effector func(context.Context) (string, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string, error) {
		var (
			response string
			err      error
		)

		//If we haven't hit our maximum retries, try again...
		for i := 0; i < retries; i++ {
			response, err = effector(ctx)
			if err != nil {
				log.Printf("Attempt %d %s; retrying in %v", i+1, err.Error(), delay)
				select {
				case <-time.After(delay):
				case <-ctx.Done():
					return "", ctx.Err()
				}
			} else {
				log.Printf("✅ Succeeded after %d attempts ✅", i+1)
				return response, nil
			}
		}

		log.Printf("Maximum attempts reached: %d failures.", retries)
		return "", err
	}
}
