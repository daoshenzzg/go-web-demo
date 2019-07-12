package dao

import (
	"context"
	"go-web-demo/library/cache/redis"
	"go-web-demo/library/log"
)

func (d *Dao) GetKey(c context.Context, key string) (val string, err error) {
	conn := d.redis.Get()
	defer conn.Close()
	if val, err = redis.String(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("GetKey(%s) error(%v)", key, err)
		}
		return
	}
	return
}

func (d *Dao) SetKey(c context.Context, key, val string, expire int64) (err error) {
	conn := d.redis.Get()
	defer conn.Close()
	if _, err = redis.String(conn.Do("SETEX", key, expire, val)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("SETEX(%s) error(%s)", key, err)
		}
		return
	}
	return
}
