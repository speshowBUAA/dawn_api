package main

import (
	"github.com/gin-gonic/gin"
	"dawn_api/router"
	"dawn_api/db"
	"dawn_api/middleware"
)

func main() {
	db.NewDawnApiDataClient()
	db.NewMongoDBClient()
	
	listener := gin.Default()

	r := listener.Group("/api")
	r.Use(jwt.JWTAuth())
	r.GET("/currentUser", router.CurrentUser)
	r.GET("/getUsers", router.GetUsers)
	r.POST("/saveAnnotation", router.SaveAnnotationCallback)

	listener.POST("/login/account", router.Login)
	listener.POST("/login/regist", router.RegisterUser)

	listener.Run(":3000")
}
