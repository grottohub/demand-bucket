package main

// import localserver "demand-bucket/local-server"
import (
	"demand-bucket/cache"
	"fmt"
)

func main() {
	// localserver.Start(8080)
	cache.AddBucket("abcdefg")
	fmt.Println(cache.GetBucket("abcdefg"))
}
