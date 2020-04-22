package optredis

import (
	"fmt"
	"github.com/sinlovgo/optredis/optredisconfig"
)

var redisConfigList *[]optredisconfig.Cfg

var errRedisConfigListEmpty = fmt.Errorf("optredis err: redis config list is empty, you must use optredis.InitByConfigList()")

// init by viper config list as
//	redis_clients:
//  - name: default
//    addr: localhost:6379
//    password:
//    db: 0
//    max_retries: 0 # Default is to not retry failed commands
//    dial_timeout: 5 # Default is 5 seconds.
//    read_timeout: 3 # Default is 3 seconds.
//    write_timeout: 3 # Default is ReadTimeout
func InitByConfigList(redisClientList []optredisconfig.Cfg) error {
	if redisClientList == nil {
		return errRedisConfigListEmpty
	}
	if len(redisClientList) == 0 {
		return errRedisConfigListEmpty
	}
	redisConfigList = &redisClientList
	return nil
}
