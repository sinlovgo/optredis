package optredis

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

const (
	envKeyCacheRedisPrefix       = "CACHE_REDIS"
	envKeyCacheRedisAddr         = "ADDR"
	envKeyCacheRedisPassword     = "PASSWORD"
	envKeyCacheRedisDB           = "DB"
	envKeyCacheRedisMaxRetries   = "MAX_RETRIES"
	envKeyCacheRedisDialTimeout  = "DIAL_TIMEOUT"
	envKeyCacheRedisReadTimeout  = "READ_TIMEOUT"
	envKeyCacheRedisWriteTimeout = "WRITE_TIMEOUT"
)

type Config struct {
	Name           string
	UseBloomFilter bool
	BloomK         uint
	BloomN         uint
	BloomM         uint
}

var defaultConfig = setDefaultConfig()

type ConfigOption func(*Config)

func setDefaultConfig() *Config {
	return &Config{
		Name:           "default",
		UseBloomFilter: false,
		BloomK:         20,
		BloomN:         1000,
		BloomM:         5,
	}
}

func NewConfig(opts ...ConfigOption) (opt *Config) {
	opt = defaultConfig
	for _, o := range opts {
		o(opt)
	}
	defaultConfig = setDefaultConfig()
	return
}

func WithName(name string) ConfigOption {
	return func(o *Config) {
		o.Name = name
	}
}

func WithUseBloomFilter(useBloomFilter bool) ConfigOption {
	return func(o *Config) {
		o.UseBloomFilter = useBloomFilter
	}
}

func WithUseBloomK(bloomK uint) ConfigOption {
	return func(o *Config) {
		o.BloomK = bloomK
	}
}

func WithUseBloomN(bloomN uint) ConfigOption {
	return func(o *Config) {
		o.BloomN = bloomN
	}
}

func WithUseBloomM(bloomM uint) ConfigOption {
	return func(o *Config) {
		o.BloomM = bloomM
	}
}

func NewConfigFull() *Config {
	return setDefaultConfig()
}

func parseEnvStringOrDefault(envKey, defaultStr string) string {
	var result string
	if viper.GetString(envKey) == "" {
		result = defaultStr
	} else {
		result = viper.GetString(envKey)
	}
	return result
}

func parseEnvIntOrDefault(envKey string, defaultInt int) int {
	var result int
	if viper.GetInt(envKey) == 0 {
		result = defaultInt
	} else {
		result = viper.GetInt(envKey)
	}
	return result
}

func parseEnvByName(name string, last string) string {
	return fmt.Sprintf("%s_%s_%s", envKeyCacheRedisPrefix, strings.ToUpper(name), last)
}
