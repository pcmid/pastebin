package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var Db *gorm.DB

func DbInit() *gorm.DB {
	db, err := gorm.Open("sqlite3", "/tmp/pastebin.db")

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

	// 设置起始id
	db.Exec("INSERT INTO 'sqlite_sequence' (`name`,`seq`) VALUES ('pastes',300000);")

	Db = db
	return db
}
