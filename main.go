package main

import (
	. "github.com/HakimovBulat/lazyCalculator/utils"
	"github.com/gin-gonic/gin"
)

var db = make(map[string]int)

func setupRouter() *gin.Engine {
	SetupLogger()
	router := gin.Default()
	router.GET("/", inputExpression)
	router.POST("/", listExpressions)
	router.LoadHTMLGlob("templates/*")
	return router
}
func main() {
	router := setupRouter()
	if err := router.Run(":8080"); err != nil {
		Logger.Error(err.Error())
	}
}

func inputExpression(c *gin.Context) {
	c.HTML(200, "index.html", nil)

}
func listExpressions(c *gin.Context) {
	math := c.PostForm("math")
	Logger.Info(math)
	c.HTML(200, "list.html", nil)
}
