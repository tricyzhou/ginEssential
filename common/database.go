package common

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"tricyzhou.com/ginessential/model"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := ""
	port := "3306"
	database := "ginessential"
	username := "root"
	password := ""
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username, password, host, port, database, charset)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	db.AutoMigrate(&model.User{}) // 自动建表
	DB = db
	return db
}

func GetDB() *gorm.DB {
	InitDB()
	return DB
}
