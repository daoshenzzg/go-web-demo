package sql

import (
	"go-web-demo/library/log"
	xtime "go-web-demo/library/time"
)

// Config mysql config.
type Config struct {
	Addr         string         // for trace
	DSN          string         // write data source name.
	ReadDSN      []string       // read data source name.
	Active       int            // pool
	Idle         int            // pool
	IdleTimeout  xtime.Duration // connect max life time.
	QueryTimeout xtime.Duration // query sql timeout
	ExecTimeout  xtime.Duration // execute sql timeout
	TranTimeout  xtime.Duration // transaction sql timeout
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *Config) (db *DB) {
	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
		panic("mysql must be set query/execute/transaction timeout")
	}
	db, err := Open(c)
	if err != nil {
		log.Error("open mysql error(%v)", err)
		panic(err)
	}
	return
}
