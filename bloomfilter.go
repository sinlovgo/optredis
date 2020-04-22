package optredis

import "github.com/willf/bloom"

// init redis filter, like of 20, 5 keys (k, m) in 1000 (n)
//	k -> the number of hashing functions on elements of the set, The actual hashing functions are important, too, but this is not a parameter for this implementation
//	n -> data size
//	m -> maximum size, typically a reasonably large multiple of the cardinality of the set to represent
func initRedisFilter(k, n, m uint) *bloom.BloomFilter {
	return bloom.New(k*n, m)
}
