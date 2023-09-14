package main

import (
	"date_calculation/controller"
	"date_calculation/middleware"

	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	serveApplication()
}

func serveApplication() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1", "23.254.209.206"})

	router.Use(middleware.CORSMiddleware())

	publicRoutes := router.Group("/api")
	publicRoutes.Use(middleware.CORSMiddleware())
	publicRoutes.POST("/CalcCalendarDate", controller.CalcCalendarDate)
	publicRoutes.POST("/CalcHundredYearDate", controller.CalcHundreYearDate)

	router.RunTLS(":8010", "./fullchain.pem", "privkey.pem")
	fmt.Println("Server running on port 8010")
}
