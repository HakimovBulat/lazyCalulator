package main

import (
	. "github.com/HakimovBulat/lazyCalculator/utils"
	"github.com/apaxa-go/eval"
	"github.com/gin-gonic/gin"
)

type Expression struct {
	Id            int
	StringVersion string
	Answer        int
	Status        string
}
type DataBase struct {
	Title       string
	Expressions []Expression
}

var id int = 1
var db DataBase
var mapOperatorsTime = map[string]int{
	"-": 100, "+": 100, "*": 100, "/": 100,
}

func setupRouter() *gin.Engine {
	SetupLogger()
	router := gin.Default()
	router.GET("/", inputExpression)
	router.POST("/", createExpression)
	router.GET("/operators", operators)
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
	c.HTML(200, "index.html", db)
}
func createExpression(c *gin.Context) {
	var newExpression Expression
	math := c.PostForm("math")
	Logger.Info(math)
	expr, err := eval.ParseString(math, "")
	if err != nil {
		newExpression.Status = "cancel"
	}
	answer, err := expr.EvalToInterface(nil)
	if err != nil {
		newExpression.Status = "cancel"
	}
	newExpression = Expression{Id: id, StringVersion: math, Answer: answer.(int), Status: "ok"}
	db.Expressions = append(db.Expressions, newExpression)
	id++
	c.HTML(200, "index.html", db)
}
func operators(c *gin.Context) {
	c.HTML(200, "operators.html", nil)
}
