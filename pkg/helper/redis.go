package helper

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func DeleteRedisKeysByPattern(client *redis.Client, pattern string) error {
	var cursor uint64
	var keys []string
	for {
		var err error
		keys, cursor, err = client.Scan(context.Background(), cursor, pattern, 10).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			_, err := client.Del(context.Background(), key).Result()
			if err != nil {
				return err
			}
			fmt.Println("Deleted key:", key)
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
