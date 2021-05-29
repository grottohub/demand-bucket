package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	fmt.Printf("\n > Connecting to Redis...\n")
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func GetBucket(key string) []string {
	reqs, err := rdb.LRange(ctx, key, 0, 19).Result()
	if err != nil {
		panic(err)
	}

	return reqs
}

func setBucketExp(key string) {
	exp, _ := time.ParseDuration("48h")
	_, err := rdb.Expire(ctx, key, exp).Result()

	if err != nil {
		fmt.Printf(" > error setting expiration %v for bucket %v\n", exp, key)
		fmt.Println(err)
	}
}

func AddBucket(key string) {
	l, err := rdb.LPush(ctx, key, "null").Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf(" > bucket %v initialized with a size of %v and initial value 'null'\n", key, l)

	setBucketExp(key)
}

func AddRequest(key string, request interface{}) {
	_, err := rdb.LPush(ctx, key, request).Result()
	if err != nil {
		panic(err)
	}

	rdb.LTrim(ctx, key, 0, 19)

	fmt.Printf(" > added request to bucket %v\n", key)
}
