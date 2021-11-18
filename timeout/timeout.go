package timeout

import (
	"context"
	"log"
)

type SlowFunction func(string) (string, error)

type WithContext func(context.Context, string) (string, error)

func Timeout(f SlowFunction) WithContext {
	return func(ctx context.Context, arg string) (string, error) {
		chres := make(chan string)
		cherr := make(chan error)

		go func() {
			res, err := f(arg)
			chres <- res
			cherr <- err
		}()

		select {
		case res := <-chres:
			log.Println("✅ Function call returned in a timely manner ✅")
			return res, <-cherr
		case <-ctx.Done():
			log.Println("❌ Function call timed out ❌")
			return "", ctx.Err()
		}
	}
}
