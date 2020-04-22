package optredis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sinlovgo/optredis/optredisconfig"
	"github.com/willf/bloom"
	"time"
)

type OptRedis struct {
	Name          string
	UseBoomFilter bool
	BloomK        uint
	BloomN        uint
	BloomM        uint

	RedisClient *redis.Client
	BloomFilter *bloom.BloomFilter

	errClientNotInit error
	errKeyEmpty      error
	errKeyFilter     error
	errKeyNotExist   error
	errDataEmpty     error
}

func (o OptRedis) Exists(key string, prefix string) (bool, error) {
	if key == "" {
		return false, o.errKeyEmpty
	}
	totalKey := fmt.Sprintf("%v%v", prefix, key)
	if o.UseBoomFilter {
		if !o.BloomFilter.TestString(totalKey) {
			return false, o.errKeyFilter
		}
	}
	return ExistsKey(o.RedisClient, totalKey)
}

func (o OptRedis) Del(key string, prefix string) (int64, error) {
	if key == "" {
		return 0, o.errKeyEmpty
	}
	totalKey := fmt.Sprintf("%v%v", prefix, key)
	return o.RedisClient.Del(totalKey).Result()
}

func (o OptRedis) TTL(key string, prefix string) (time.Duration, error) {
	if key == "" {
		return 0, o.errKeyEmpty
	}
	totalKey := fmt.Sprintf("%v%v", prefix, key)
	if o.UseBoomFilter {
		if !o.BloomFilter.TestString(totalKey) {
			return 0, o.errKeyFilter
		}
	}
	existCount, err := o.RedisClient.Exists(totalKey).Result()
	if err != nil {
		return 0, err
	}
	if existCount == 0 {
		return 0, o.errKeyNotExist
	}
	return o.RedisClient.TTL(totalKey).Result()
}

func (o OptRedis) Expire(key string, prefix string, exp time.Duration) (bool, error) {
	if key == "" {
		return false, o.errKeyEmpty
	}
	totalKey := fmt.Sprintf("%v%v", prefix, key)
	if o.UseBoomFilter {
		if !o.BloomFilter.TestString(totalKey) {
			return false, o.errKeyFilter
		}
	}
	existCount, err := o.RedisClient.Exists(totalKey).Result()
	if err != nil {
		return false, err
	}
	if existCount == 0 {
		return false, o.errKeyNotExist
	}
	return o.RedisClient.Expire(totalKey, exp).Result()
}

func (o OptRedis) Persist(key string, prefix string) (bool, error) {
	if key == "" {
		return false, o.errKeyEmpty
	}
	totalKey := fmt.Sprintf("%v%v", prefix, key)
	if o.UseBoomFilter {
		if !o.BloomFilter.TestString(totalKey) {
			return false, o.errKeyFilter
		}
	}
	return o.RedisClient.Persist(key).Result()
}

func (o OptRedis) SetJson(key string, prefix string, data interface{}, expiration time.Duration) error {
	if key == "" {
		return o.errKeyEmpty
	}
	if data == nil {
		return o.errDataEmpty
	}
	if o.RedisClient == nil {
		return o.errClientNotInit
	}

	marshal, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	totalKey := fmt.Sprintf(`%v%v`, prefix, key)
	err = o.RedisClient.
		Set(totalKey, string(marshal), expiration).
		Err()
	if err != nil {
		return err
	}
	if o.UseBoomFilter {
		o.BloomFilter.AddString(totalKey)
	}
	return nil
}

func (o OptRedis) GetJson(key string, prefix string, v interface{}) error {
	if key == "" {
		return o.errKeyEmpty
	}
	if o.RedisClient == nil {
		return o.errClientNotInit
	}
	totalKey := fmt.Sprintf("%v%v", prefix, key)
	if o.UseBoomFilter {
		if !o.BloomFilter.TestString(totalKey) {
			return o.errKeyFilter
		}
	}
	result, err := o.RedisClient.Get(totalKey).Result()
	if err != nil {
		return err
	}
	if result == "" {
		return o.errDataEmpty
	}
	err = json.Unmarshal([]byte(result), &v)
	if err != nil {
		return err
	}
	return nil
}

func (o OptRedis) Client() *redis.Client {
	return o.RedisClient
}

func (o OptRedis) Ping() (OptRedis, error) {
	_, err := o.RedisClient.Ping().Result()
	if err != nil {
		return o, err
	}
	return o, nil
}

func (o OptRedis) InitByName() OptRedis {
	if redisConfigList == nil {
		panic(errRedisConfigListEmpty)
	}
	redisConf := optredisconfig.ByName(*redisConfigList, o.Name)
	redisAddr := parseEnvStringOrDefault(parseEnvByName(o.Name, envKeyCacheRedisAddr), redisConf.Addr)
	redisPassword := parseEnvStringOrDefault(parseEnvByName(o.Name, envKeyCacheRedisPassword), redisConf.Password)
	redisDB := parseEnvIntOrDefault(parseEnvByName(o.Name, envKeyCacheRedisDB), redisConf.DB)
	redisMaxRetries := parseEnvIntOrDefault(parseEnvByName(o.Name, envKeyCacheRedisMaxRetries), redisConf.MaxRetries)
	redisDialTimeout := parseEnvIntOrDefault(parseEnvByName(o.Name, envKeyCacheRedisDialTimeout), redisConf.DialTimeout)
	redisReadTimeout := parseEnvIntOrDefault(parseEnvByName(o.Name, envKeyCacheRedisReadTimeout), redisConf.ReadTimeout)
	redisWriteTimeout := parseEnvIntOrDefault(parseEnvByName(o.Name, envKeyCacheRedisWriteTimeout), redisConf.WriteTimeout)

	o.RedisClient = redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     redisPassword, // no password set ""
		DB:           redisDB,       // use default DB 0
		MaxRetries:   redisMaxRetries,
		DialTimeout:  time.Duration(redisDialTimeout) * time.Second,
		ReadTimeout:  time.Duration(redisReadTimeout) * time.Second,
		WriteTimeout: time.Duration(redisWriteTimeout) * time.Second,
	})

	if o.UseBoomFilter {
		o.BloomFilter = initRedisFilter(5, 100, 20)
	}

	return o
}

type RedisFunc interface {
	InitByName() OptRedis
	Ping() (OptRedis, error)
	Client() *redis.Client
	Exists(key string, prefix string) (bool, error)
	Del(key string, prefix string) (int64, error)
	TTL(key string, prefix string) (time.Duration, error)
	Expire(key string, prefix string, exp time.Duration) (bool, error)
	Persist(key string, prefix string) (bool, error)
	SetJson(key string, prefix string, data interface{}, expiration time.Duration) error
	GetJson(key string, prefix string, v interface{}) error
}

func NewOptRedis(cfg Config) RedisFunc {
	return &OptRedis{
		Name:          cfg.Name,
		UseBoomFilter: cfg.UseBloomFilter,

		errClientNotInit: fmt.Errorf("%v opt redis : client empty, must be init", cfg.Name),
		errKeyEmpty:      fmt.Errorf("%v opt redis : key is empty plase check", cfg.Name),
		errKeyFilter:     fmt.Errorf("%v opt redis : key filtered", cfg.Name),
		errKeyNotExist:   fmt.Errorf("%v opt redis : key not exist", cfg.Name),
		errDataEmpty:     fmt.Errorf("%v opt redis : data is empty", cfg.Name),
	}
}
