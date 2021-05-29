package cache

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrcSB(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())

	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

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

func AddBucket() string {
	key := randStringBytesMaskImprSrcSB(8)

	l, err := rdb.LPush(ctx, key, "null").Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf(" > bucket %v initialized with a size of %v and initial value 'null'\n", key, l)

	setBucketExp(key)

	return key
}

func AddRequest(key string, request interface{}) {
	_, err := rdb.LPush(ctx, key, request).Result()
	if err != nil {
		panic(err)
	}

	rdb.LTrim(ctx, key, 0, 19)

	fmt.Printf(" > added request to bucket %v\n", key)
}
