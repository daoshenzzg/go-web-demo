package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"person.mgtv.com/framework/config"
	"person.mgtv.com/framework/logs"
)

var mapping = make(map[string]*sql.DB)

func init() {
	// init feed db
	InitDatabase("db.feed")
}

func InitDatabase(sectionKey string) {
	section := config.Section(sectionKey)

	address := section.Key("address").String()
	maxIdel, _ := section.Key("max_idel").Int()
	maxConn, _ := section.Key("max_conn").Int()

	db, err := Open(address, maxIdel, maxConn)
	if err != nil {
		panic(err)
	}
	mapping[sectionKey] = db
	logs.GetLogger("system").Infof("Database[%s] init success.", address)
}

func Open(address string, maxOpenConns, maxIdleConns int) (*sql.DB, error) {
	db, err := sql.Open("mysql", address)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Get(name string) *sql.DB {
	return mapping[name]
}

func String(v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	return v.String
}

func Int(v sql.NullInt64) int {
	if !v.Valid {
		return 0
	}
	return int(v.Int64)
}

func Int64(v sql.NullInt64) int64 {
	if !v.Valid {
		return 0
	}
	return v.Int64
}

func Float64(v sql.NullFloat64) float64 {
	if !v.Valid {
		return 0
	}
	return v.Float64
}
