package main

import "github.com/pcmid/pasterbin/model"

func main() {
	db := model.DbInit()
	defer db.Close()
	_ = routerInit().Run()
}
