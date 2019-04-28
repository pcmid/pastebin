package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pcmid/pastebin/controller"
	"github.com/pcmid/pastebin/model"
)

func routerInit() *gin.Engine {
	r := gin.Default()

	r.POST("/",controller.ContentTypeCheck, model.CreatePaste)

	r.GET("/:hash", model.FetchPaste)

	r.PUT("/:hash",controller.ContentTypeCheck, model.UpdatePaste)

	r.DELETE("/:hash", model.DeletePaste)

	return r

}
