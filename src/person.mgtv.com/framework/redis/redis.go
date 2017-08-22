package redis

import (
	"strings"
	"time"

	"github.com/go-redis/redis"
	"person.mgtv.com/framework/config"
	"person.mgtv.com/framework/logs"
)

var mapping = make(map[string]*redis.ClusterClient)

func init() {
	// init feed redis cluster
	InitRedisCluster("redis.feed")
}

func InitRedisCluster(sectionKey string) {
	section := config.Section(sectionKey)

	address := section.Key("address").String()
	timeout, _ := section.Key("timeout").Duration()
	poolSize, _ := section.Key("pool_size").Int()

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       strings.Split(address, ","),
		MaxRetries:  3,
		DialTimeout: 1 * time.Duration(time.Second),
		ReadTimeout: timeout * time.Duration(time.Millisecond),
		PoolSize:    poolSize,
	})

	mapping[sectionKey] = client
	logs.GetLogger("system").Infof("RedisCluster[%s] [%s] init success.", sectionKey, address)
}

func GetRedis(name string) *redis.ClusterClient {
	if val, ok := mapping[name]; ok {
		return val
	}

	logs.GetLogger("system").Errorf("Redis[%s] is not exists.", name)
	return nil
}
