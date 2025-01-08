package db

import (
	"github.com/garyburd/redigo/redis"
	"github.com/qingw1230/studyim/pkg/common/constant"
	"github.com/qingw1230/studyim/pkg/common/log"
)

const (
	uidPidToken = "UID_PID_TOKEN_STATUS:"
)

func (d *DataBases) Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	conn := d.redisPool.Get()
	if err := conn.Err(); err != nil {
		log.Error("", "redis cmd = %v, err = %v", cmd, err)
		return nil, err
	}
	defer conn.Close()

	params := make([]interface{}, 0)
	params = append(params, key)

	if len(args) > 0 {
		params = append(params, args...)
	}
	return conn.Do(cmd, params...)
}

// AddTokenFlag store userID and flatform class to redis.
func (d *DataBases) AddTokenFlag(userID string, platformID int32, token string, flag int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	log.Debug("", "add token key is: %s", key)
	_, err := d.Exec("HSET", key, token, flag)
	return err
}

func (d *DataBases) GetTokenMapByUidPid(userID, platformID string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platformID
	return redis.IntMap(d.Exec("HGETALL", key))
}

func (d *DataBases) SetTokenMapByUidPid(userID string, platformID int32, m map[string]int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	_, err := d.Exec("HMSET", key, redis.Args{}.Add().AddFlat(m)...)
	return err
}

func (d *DataBases) DeleteTokenByUidPid(userID string, platformID int32, fields []string) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	_, err := d.Exec("HDEL", key, redis.Args{}.Add().AddFlat(fields)...)
	return err
}
