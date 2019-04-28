package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Db *gorm.DB

func DbInit() *gorm.DB {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=pastebin password=s sslmode=disable")

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(true)

	/*
	 * @date: 19-4-27 下午7:25
	 * @user: id
	 * @todo: 添加AutoMigrate()
	 */

	db.AutoMigrate(Paste{})

	Db = db
	return db
}
