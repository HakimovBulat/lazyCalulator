package router

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/HakimovBulat/lazyCalculator/utils"
	"github.com/apaxa-go/eval"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func SetupRouter() *gin.Engine {
	var err error
	utils.SetupLogger()
	connection, err = sql.Open("postgres", connectionString)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	router := gin.Default()
	router.GET("/", inputExpression)
	router.POST("/", createExpression)
	router.GET("/operators", showOperatorsTime)
	router.PUT("/operators", replaceOperatorsTime)
	router.GET("/static_operators", operatorsStatic)
	router.POST("/static_operators", operatorsStatic)
	router.GET("/get_expression/:id", getExpression)
	router.LoadHTMLGlob("templates/*.html")
	return router
}

type Expression struct {
	Id            int
	StringVersion string
	Answer        string
	Status        string
	StartDate     time.Time
	EndDate       time.Time
}

var mapOperatorsTime = map[string]int{
	"-": 10,
	"+": 10,
	"*": 10,
	"/": 10,
}
var connectionString = "host=127.0.0.1 port=5432 user=postgres password=Love_and_elephant42 dbname=Expression sslmode=disable"
var connection *sql.DB

func inputExpression(c *gin.Context) {
	rows, err := connection.Query(`SELECT * FROM "Expression"`)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	defer rows.Close()
	expressions := []Expression{}
	for rows.Next() {
		validExpression := new(Expression)
		rows.Scan(
			&validExpression.Id,
			&validExpression.StringVersion,
			&validExpression.Status,
			&validExpression.Answer,
			&validExpression.StartDate,
			&validExpression.EndDate,
		)
		now := time.Now()
		if validExpression.EndDate.Before(now) {
			if validExpression.Answer != "not found" {
				validExpression.Status = "ok"
			} else {
				validExpression.Status = "cancel"
			}
		}
		_, err = connection.Query(`UPDATE "Expression" SET "StringVersion"=$2, "Status"=$3, "Answer"=$4, "StartDate"=$5, "EndDate"=$6
			WHERE "id"=$1`,
			validExpression.Id,
			validExpression.StringVersion,
			validExpression.Status,
			validExpression.Answer,
			validExpression.StartDate,
			validExpression.EndDate,
		)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
		expressions = append(expressions, *validExpression)
	}
	c.HTML(200, "index.html", expressions)
}
func createExpression(c *gin.Context) {
	now := time.Now()
	var newExpression Expression
	mathExpression := c.PostForm("math")
	utils.Logger.Info(mathExpression)
	endDate := getTime(mathExpression, now)
	newExpression = Expression{
		Id:            0,
		StringVersion: mathExpression,
		Status:        "ok",
		StartDate:     now,
		EndDate:       endDate,
	}
	expr, err := eval.ParseString(newExpression.StringVersion, "")
	if err != nil {
		newExpression.Status = "cancel"
	} else {
		answer, err := expr.EvalToInterface(nil)
		if err != nil || answer == nil {
			newExpression.Answer = "not found"
		} else {
			newExpression.Answer = fmt.Sprint(answer)
		}
	}
	newExpression.Status = "process"
	expressions := []Expression{}
	if mathExpression != "" {
		_, err := connection.Query(`INSERT INTO "Expression"("StringVersion", "Status", "Answer", "StartDate", "EndDate")
		VALUES($1, $2, $3, $4, $5)`,
			newExpression.StringVersion,
			newExpression.Status,
			newExpression.Answer,
			newExpression.StartDate,
			newExpression.EndDate,
		)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
		rows, err := connection.Query(`SELECT * FROM "Expression"`)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
		for rows.Next() {
			validExpression := new(Expression)
			rows.Scan(
				&validExpression.Id,
				&validExpression.StringVersion,
				&validExpression.Status,
				&validExpression.Answer,
				&validExpression.StartDate,
				&validExpression.EndDate,
			)
			expressions = append(expressions, *validExpression)
		}
		rows.Close()
	}
	c.HTML(200, "index.html", expressions)
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
func getExpression(c *gin.Context) {
	var expression Expression
	connection.QueryRow(`SELECT * FROM "Expression" WHERE "id"=$1`, c.Param("id")).Scan(
		&expression.Id,
		&expression.StringVersion,
		&expression.Status,
		&expression.Answer,
		&expression.StartDate,
		&expression.EndDate,
	)
	var answer any
	var err error
	expr, err := eval.ParseString(expression.StringVersion, "")
	if err != nil {
		expression.Status = "cancel"
		c.JSON(400, expression)
		return
	}
	now := time.Now()
	if expression.EndDate.Before(now) {
		answer, err = expr.EvalToInterface(nil)
		if err != nil || answer == nil {
			expression.Status = "cancel"
			c.JSON(400, expression)
			return
		}
		expression.Answer = fmt.Sprint(answer)
		expression.Status = "ok"
	}
	c.JSON(200, expression)
}
func getTime(expression string, now time.Time) time.Time {
	var seconds int
	for _, symbol := range expression {
		if string(symbol) == "-" || string(symbol) == "+" || string(symbol) == "/" || string(symbol) == "*" {
			seconds += mapOperatorsTime[string(symbol)]
		}
	}
	return now.Add(time.Duration(seconds) * time.Second)
}
