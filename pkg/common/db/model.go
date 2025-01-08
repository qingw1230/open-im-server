package db

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/qingw1230/studyim/pkg/common/config"
)

var DB DataBases

type DataBases struct {
	MysqlDB   mysqlDB
	redisPool *redis.Pool
}

func key(dbAddress, dbName string) string {
	return dbAddress + "_" + dbName
}

func init() {
	initMysqlDB()

	DB.redisPool = &redis.Pool{
		MaxIdle:     config.Config.Redis.DBMaxIdle,
		MaxActive:   config.Config.Redis.DBMaxActive,
		IdleTimeout: time.Duration(config.Config.Redis.DBIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				config.Config.Redis.DBAddress,
				redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialDatabase(0),
				redis.DialPassword(config.Config.Redis.DBPassWord),
			)
		},
	}
}
