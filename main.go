package main

import (
	"fmt"
	"strconv"

	. "github.com/HakimovBulat/lazyCalculator/utils"
	"github.com/apaxa-go/eval"
	"github.com/gin-gonic/gin"
)

type Expression struct {
	Id            int
	StringVersion string
	Answer        string
	Status        string
}
type DataBase struct {
	Title       string
	Expressions []Expression
}

var id int = 1
var db DataBase
var mapOperatorsTime = map[string]int{
	"-": 100,
	"+": 100,
	"*": 100,
	"/": 100,
}

func setupRouter() *gin.Engine {
	SetupLogger()
	router := gin.Default()
	router.GET("/", inputExpression)
	router.POST("/", createExpression)
	router.GET("/operators", showOperatorsTime)
	router.PUT("/operators", replaceOperatorsTime)
	router.GET("/static_operators", operatorsStatic)
	router.POST("/static_operators", operatorsStatic)
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
	var answer any
	var err error
	mathExpression := c.PostForm("math")
	Logger.Info(mathExpression)
	expr, err := eval.ParseString(mathExpression, "")
	if err != nil {
		newExpression = Expression{Id: id, StringVersion: mathExpression, Answer: fmt.Sprint(answer), Status: "cancel"}
		db.Expressions = append(db.Expressions, newExpression)
		id++
		c.HTML(200, "index.html", db)
		return
	}
	answer, err = expr.EvalToInterface(nil)
	if err != nil || answer == nil {
		newExpression = Expression{Id: id, StringVersion: mathExpression, Answer: "not found", Status: "cancel"}
		db.Expressions = append(db.Expressions, newExpression)
		id++
		c.HTML(200, "index.html", db)
		return
	}
	fmt.Println(answer)
	newExpression = Expression{Id: id, StringVersion: mathExpression, Answer: fmt.Sprint(answer), Status: "ok"}
	db.Expressions = append(db.Expressions, newExpression)
	id++
	c.HTML(200, "index.html", db)
}
func showOperatorsTime(c *gin.Context) {
	c.HTML(200, "operators.html", gin.H{
		"addition":       mapOperatorsTime["+"],
		"substraction":   mapOperatorsTime["-"],
		"multiplication": mapOperatorsTime["*"],
		"division":       mapOperatorsTime["/"],
	})
}

func replaceOperatorsTime(c *gin.Context) {
	mapOperatorsTime["+"], _ = strconv.Atoi(c.PostForm("addition"))
	mapOperatorsTime["-"], _ = strconv.Atoi(c.PostForm("substraction"))
	mapOperatorsTime["*"], _ = strconv.Atoi(c.PostForm("multiplication"))
	mapOperatorsTime["/"], _ = strconv.Atoi(c.PostForm("division"))
	c.HTML(200, "operators.html", gin.H{
		"addition":       mapOperatorsTime["+"],
		"substraction":   mapOperatorsTime["-"],
		"multiplication": mapOperatorsTime["*"],
		"division":       mapOperatorsTime["/"],
	})
}
func operatorsStatic(c *gin.Context) {
	if c.Request.Method == "POST" {
		mapOperatorsTime["+"], _ = strconv.Atoi(c.PostForm("addition"))
		mapOperatorsTime["-"], _ = strconv.Atoi(c.PostForm("substraction"))
		mapOperatorsTime["*"], _ = strconv.Atoi(c.PostForm("multiplication"))
		mapOperatorsTime["/"], _ = strconv.Atoi(c.PostForm("division"))
	}
	c.HTML(200, "static_operators.html", gin.H{
		"addition":       mapOperatorsTime["+"],
		"substraction":   mapOperatorsTime["-"],
		"multiplication": mapOperatorsTime["*"],
		"division":       mapOperatorsTime["/"],
	})
}
