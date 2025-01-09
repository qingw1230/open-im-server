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
	log.Info("open db ok ", dsn)

	db.AutoMigrate(
		&User{},
	)
	db.Set("gorm:table_options", "CHARSET=utf8")
	db.Set("gorm:table_options", "collation=utf8_unicode_ci")

	if !db.HasTable(&User{}) {
		log.Info("CreateTable User")
		db.CreateTable(&User{})
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
