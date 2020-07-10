package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"

	. "github.com/zzsds/micro-upload/conf"
	"github.com/zzsds/micro-upload/handler"
)

func main() {
	InitConfig()

	// New Service
	service := web.NewService(
		web.Name("go.micro.api.upload"),
		web.Version("latest"),
	)

	service.Init()
	// Create RESTful handler (using Gin)

	up := &handler.Upload{}
	router := gin.Default()
	router.POST("/upload/aliyun", up.Aliyun)

	// Register Handler
	service.Handle("/", router)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
