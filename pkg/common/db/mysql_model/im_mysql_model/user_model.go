package im_mysql_model

import (
	"time"

	"github.com/qingw1230/studyim/pkg/common/db"
)

type User struct {
	UserID      string    `gorm:"column:user_id;primaryKey;"`
	Nickname    string    `gorm:"column:name"`
	FaceUrl     string    `gorm:"column:icon"`
	Gender      int32     `gorm:"column:gender"`
	PhoneNumber string    `gorm:"column:phone_number"`
	Birth       string    `gorm:"column:birth"`
	Email       string    `gorm:"column:email"`
	Ex          string    `gorm:"column:ex"`
	CreateTime  time.Time `gorm:"column:create_time"`
}

func GetUserByUserID(userID string) (*db.User, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var user db.User
	err = dbConn.Table("users").Where("user_id=?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
