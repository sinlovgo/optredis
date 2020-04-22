[![TravisBuildStatus](https://api.travis-ci.org/sinlovgo/optredis.svg?branch=master)](https://travis-ci.org/sinlovgo/optredis)
[![GoDoc](https://godoc.org/github.com/sinlovgo/optredis?status.png)](https://godoc.org/github.com/sinlovgo/optredis/)
[![GoReportCard](https://goreportcard.com/badge/github.com/sinlovgo/optredis)](https://goreportcard.com/report/github.com/sinlovgo/optredis)
[![codecov](https://codecov.io/gh/sinlovgo/optredis/branch/master/graph/badge.svg)](https://codecov.io/gh/sinlovgo/optredis)

## for what

- redis client loader and support bloom filter

## depends

in go mod project

```bash
# warning use privte git host must set
# global set for once
# add private git host like github.com to evn GOPRIVATE
$ go env -w GOPRIVATE='github.com'
# use ssh proxy
# set ssh-key to use ssh as http
$ git config --global url."git@github.com:".insteadOf "http://github.com/"
# or use PRIVATE-TOKEN
# set PRIVATE-TOKEN as gitlab or gitea
$ git config --global http.extraheader "PRIVATE-TOKEN: {PRIVATE-TOKEN}"
# set this rep to download ssh as https use PRIVATE-TOKEN
$ git config --global url."ssh://github.com/".insteadOf "https://github.com/"

# before above global settings
# test version info
$ git ls-remote -q http://github.com/sinlovgo/optredis.git

# test depends see full version
$ go list -v -m -versions github.com/sinlovgo/optredis
# or use last version add go.mod by script
$ echo "go mod edit -require=$(go list -m -versions github.com/sinlovgo/optredis | awk '{print $1 "@" $NF}')"
$ echo "go mod vendor"
```

## evn

- golang sdk 1.13+
- github.com/go-redis/redis v6.15.7
- github.com/spf13/viper v1.6.3
- github.com/willf/bitset v1.1.10
- github.com/willf/bloom v2.0.3

## use

- optredis client and tools of `default.go`

```go
package cacheDefault
import "github.com/sinlovgo/optredis"
const (
	name string = "default"
)

var defaultOptRedis *optredis.OptRedis

func Init() error {
	if defaultOptRedis == nil {
		config := optredis.NewConfig(
			optredis.WithName("default"),
			optredis.WithUseBoomFilter(false),
			optredis.WithUseBloomK(20),
			optredis.WithUseBloomN(1000),
			optredis.WithUseBloomM(5),
		)
		optRedis, err := optredis.NewOptRedis(*config).InitByName().Ping()
		if err != nil {
			return err
		}
		defaultOptRedis = &optRedis
	}
	return nil
}

func Opt() *optredis.OptRedis {
	return defaultOptRedis
}
```

- init at `main.go`

```go
package main

import (
	"cacheDefault"
	"fmt"
	"github.com/sinlovgo/optredis"
)
func main()  {
	err := InitRedisOpt()
	if err!= nil {
		fmt.Printf("init optredis err: %v", err)
	}
}

func InitRedisOpt() error {
	err := optredis.InitByConfigList(cfg.Global().RedisOptConfig)
	if err != nil {
		return err
	}
	return cacheDefault.Init()
}
```

- then can use topic of cache student, student struct at package at `demo`

```go
package cacheTopicDemo
import (
	"cacheDefault"
	"demo"
	"time"
)

const (
	cpStudentPrefix string = "cache-student-"
	// one week
	cpStudentExpiration = time.Duration(24) * time.Hour
)

func ExistsStudent(key string) (bool, error) {
	return cacheDefault.Opt().Exists(key, cpStudentPrefix)
}

func Set(key string, data *demo.Student) error {
	return cacheDefault.Opt().SetJson(key, cpStudentPrefix, data, cpStudentExpiration)
}

func Get(key string, data *demo.Student) error {
	return cacheDefault.Opt().GetJson(key, cpStudentPrefix, data)
}
```

## dev

```bash
make init
```

- test code

```bash
make test
```

add main.go file and run

```bash
make run
```