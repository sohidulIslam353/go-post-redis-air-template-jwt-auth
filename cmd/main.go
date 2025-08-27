package main

import (
	"gin-app/config"
	"gin-app/internal/pkg/router"
)

func main() {
	config.InitDB()       // DB return করছে
	config.ConnectRedis() // redis connection
	r := router.SetupRouter()
	r.Run(":8080")
}
