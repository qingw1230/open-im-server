package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/qingw1230/studyim/pkg/common/config"
	"github.com/qingw1230/studyim/pkg/common/log"
)

type mysqlDB struct {
	rw    sync.RWMutex
	dbMap map[string]*gorm.DB
}

func initMysqlDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], "mysql")
	var db *gorm.DB
	var err1 error
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Error("0", "Open failed ", err.Error(), dsn)
	}
	if err != nil {
		time.Sleep(time.Duration(10) * time.Second)
		db, err1 = gorm.Open("mysql", dsn)
		if err1 != nil {
			log.Error("0", "Open failed ", err1.Error(), dsn)
			panic(err1.Error())
		}
	}

	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8 COLLATE utf8_general_ci;", config.Config.Mysql.DBDatabaseName)
	err = db.Exec(sql).Error
	if err != nil {
		log.Error("0", "Exec failed ", err.Error(), sql)
		panic(err.Error())
	}
	db.Close()

	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], config.Config.Mysql.DBDatabaseName)
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Error("0", "Open failed ", err.Error(), dsn)
		panic(err.Error())
	}

	sqlTable := "CREATE TABLE IF NOT EXISTS `user` (" +
		" `uid` varchar(64) NOT NULL," +
		" `name` varchar(64) DEFAULT NULL," +
		" `icon` varchar(1024) DEFAULT NULL," +
		" `gender` tinyint(4) unsigned zerofill DEFAULT NULL," +
		" `mobile` varchar(32) DEFAULT NULL," +
		" `birth` varchar(16) DEFAULT NULL," +
		" `email` varchar(64) DEFAULT NULL," +
		" `ex` varchar(1024) DEFAULT NULL," +
		" `create_time` datetime DEFAULT NULL," +
		" PRIMARY KEY (`uid`)," +
		" UNIQUE KEY `uk_uid` (`uid`)" +
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"
	err = db.Exec(sqlTable).Error
	if err != nil {
		panic(err.Error())
	}
}

func (m *mysqlDB) DefaultGormDB() (*gorm.DB, error) {
	return m.GormDB(config.Config.Mysql.DBAddress[0], config.Config.Mysql.DBDatabaseName)
}

func (m *mysqlDB) GormDB(dbAddress, dbName string) (*gorm.DB, error) {
	m.rw.Lock()
	defer m.rw.Unlock()

	k := key(dbAddress, dbName)
	if _, ok := m.dbMap[k]; !ok {
		if err := m.open(dbAddress, dbName); err != nil {
			return nil, err
		}
	}
	return m.dbMap[k], nil
}

func (m *mysqlDB) open(dbAddress, dbName string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, dbAddress, dbName)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}

	db.SingularTable(true)
	db.DB().SetMaxOpenConns(config.Config.Mysql.DBMaxOpenConns)
	db.DB().SetMaxIdleConns(config.Config.Mysql.DBMaxIdleConns)
	db.DB().SetConnMaxLifetime(time.Duration(config.Config.Mysql.DBMaxLifeTime) * time.Second)

	if m.dbMap == nil {
		m.dbMap = make(map[string]*gorm.DB)
	}
	k := key(dbAddress, dbName)
	m.dbMap[k] = db
	return nil
}
