package main

import (
	"github.com/gin-gonic/gin"
	"dawn_api/router"
	"dawn_api/db"
	"fmt"
)

func main() {
	db.NewDawnApiDataClient()
	fmt.Println(db.Query("id", "1"))
	
	listern := gin.Default()
	listern.POST("/api/login/account", router.TrajectoryCallback)
	listern.Run(":3000")
}
