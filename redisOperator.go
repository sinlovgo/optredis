package optredis

import (
	"fmt"
	"github.com/go-redis/redis"
)

var errRedisClientEmpty = fmt.Errorf("redist clinet is empty")

func ExistsKey(redisCli *redis.Client, key string) (bool, error) {
	if redisCli == nil {
		return false, errRedisClientEmpty
	}
	count, err := redisCli.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

// scan keys instead redisCli.Keys()
//	redisCli *redis.Client
//	match string
//	maxCount int64
// return
//	error scan error
//	[]string removed repeated key
func RedisScanKeysMatch(redisCli *redis.Client, match string, maxCount int64) ([]string, error) {
	if redisCli == nil {
		return nil, errRedisClientEmpty
	}
	var cursor uint64
	var scanFull []string
	for {
		keys, cursor, err := redisCli.Scan(cursor, match, maxCount).Result()
		if err != nil {
			return nil, err
		}
		if len(keys) > 0 {
			for _, v := range keys {
				scanFull = append(scanFull, v)
			}
		}
		if cursor == 0 {
			break
		}
	}
	scanRes := removeRepeatedElementString(scanFull)
	return scanRes, nil
}

func removeRepeatedElementString(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
