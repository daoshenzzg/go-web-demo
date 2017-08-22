package demo

import (
	"github.com/go-redis/redis"
	demoDao "person.mgtv.com/dao/demo"
	redisCluster "person.mgtv.com/framework/redis"
	demoModel "person.mgtv.com/model/demo"
	feedServer "person.mgtv.com/thirdparty/feed"
	"person.mgtv.com/framework/logs"
	"time"
)

type DemoService struct {
	demoDao   *demoDao.DemoDao
	feedRedis *redis.ClusterClient
}

func NewDemoService() *DemoService {
	return &DemoService{
		demoDao:   demoDao.NewDemoDao(),
		feedRedis: redisCluster.GetRedis("redis.feed"),
	}
}

// db and redis  demo
func (service *DemoService) GetFeed(feedId string) (feed *demoModel.MaxTimeline, err error) {
	feed, err = service.demoDao.GetFeed(feedId)
	if err != nil {
		return nil, err
	}

	return feed, err
}

func (service *DemoService) MutiCommit() error {
	// 开启事务1
	tx, err := service.demoDao.DB.Begin()
	if err != nil {
		return err
	}

	// commit step 1
	err = service.demoDao.UpdateSeq(tx, 36, "11111111")
	if err != nil {
		tx.Rollback()
		return err
	}

	// commit step 2
	err = service.demoDao.UpdateSeq(tx, 37, "22222222")
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	logs.GetLogger("system").Debugf("Multi commit success.")

	return nil
}

// redis demo
func (service *DemoService) GetKey(key string) (string, error) {
	val := service.feedRedis.Get("test_key").Val()
	if val == "" {
		service.feedRedis.Set("test_key", "你好你好，我是从Redis中查出来的", time.Duration(time.Minute))
		val = service.feedRedis.Get("test_key").Val()
	}

	logs.GetLogger("system").Debug("Redis test:", val)
	return val, nil
}

// httpClient demo
func (service *DemoService) IsFollowed(uid, artistId string) (bool, error) {
	isFollowed, err := feedServer.IsFollowed(uid, artistId)
	if err != nil {
		return false, err
	}
	return isFollowed, err
}
