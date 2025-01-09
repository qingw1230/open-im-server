package im_mysql_model

import (
	"time"

	"github.com/qingw1230/studyim/pkg/common/constant"
	"github.com/qingw1230/studyim/pkg/common/db"
	"github.com/qingw1230/studyim/pkg/utils"
)

// type User struct {
// 	UserID      string    `gorm:"column:user_id;primaryKey;"`
// 	Nickname    string    `gorm:"column:name"`
// 	FaceUrl     string    `gorm:"column:icon"`
// 	Gender      int32     `gorm:"column:gender"`
// 	PhoneNumber string    `gorm:"column:phone_number"`
// 	Birth       string    `gorm:"column:birth"`
// 	Email       string    `gorm:"column:email"`
// 	Ex          string    `gorm:"column:ex"`
// 	CreateTime  time.Time `gorm:"column:create_time"`
// }

func UserRegister(user db.User) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	user.CreateTime = time.Now()
	if user.AppMangerLevel == 0 {
		user.AppMangerLevel = constant.AppOrdinaryUsers
	}
	if user.Birth.Unix() < 0 {
		user.Birth = utils.UnixSecondToTime(0)
	}
	err = dbConn.Table("users").Create(&user).Error
	return err
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
