package dao

import (
	"context"
	"go-web-demo/app/demo-api/internal/conf"
	"go-web-demo/library/cache/redis"
	xsql "go-web-demo/library/database/sql"
	xhttp "go-web-demo/library/http"
	"go-web-demo/library/log"
	"net/http"
)

// Dao struct
type Dao struct {
	c *conf.Config
	// mysql
	db *xsql.DB
	// redis
	redis *redis.Pool
	// httpClient
	httpClient *http.Client
}

// New init
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:          c,
		db:         xsql.NewMySQL(c.MySQL.School),
		redis:      redis.NewPool(c.Redis),
		httpClient: xhttp.NewClient(c.HttpClient.Paopao),
	}
	return
}

//BeginTran begin transaction
func (d *Dao) BeginTran(ctx context.Context) (tx *xsql.Tx, err error) {
	if tx, err = d.db.Begin(ctx); err != nil {
		log.Error("BeginTran d.arcDB.Begin error(%v)", err)
	}
	return
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) (err error) {
	// TODO
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.db.Close()
}
