package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pcmid/pastebin/controller"
	"github.com/pcmid/pastebin/model"
)

func routerInit() *gin.Engine {
	r := gin.Default()

	r.POST("/", controller.ContentTypeCheck, model.CreatePaste)

	r.GET("/:url", model.FetchPaste)

	r.GET("/:url/meta", model.FetchMeta)

	r.PUT("/:url", controller.ContentTypeCheck, model.UpdatePaste)

	r.DELETE("/:url", model.DeletePaste)

	return r

}
