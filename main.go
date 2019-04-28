package main

import "github.com/pcmid/pastebin/model"

func main() {
	db := model.DbInit()
	defer db.Close()
	_ = routerInit().Run()
}
