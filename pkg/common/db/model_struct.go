package db

import "time"

type User struct {
	UserID         string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname       string    `gorm:"column:name;size:255"`
	FaceURL        string    `gorm:"column:face_url;size:255"`
	Gender         int32     `gorm:"column:gender"`
	PhoneNumber    string    `gorm:"column:phone_number;size:32"`
	Birth          time.Time `gorm:"column:birth"`
	Email          string    `gorm:"column:email;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AppMangerLevel int32     `gorm:"column:app_manger_level"`
}
