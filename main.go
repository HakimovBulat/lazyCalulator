package main

import (
	"fmt"
	"net/http"

	"github.com/apaxa-go/eval"
	"go.uber.org/zap"
)

func indexHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "POST" {
		expression := r.FormValue("name")
		logger, err := zap.NewDevelopment()
		if err != nil {
			fmt.Println("wrong logger")
		}
		expressionEval, err := eval.ParseString(expression, "")
		if err != nil {
			logger.Error(err.Error())
		}
		answer, err := expressionEval.EvalToInterface(nil)
		logger.Info(fmt.Sprint(answer), zap.String("expression", expression))
	}
}
func main() {
	http.HandleFunc("/", indexHandle)
	http.ListenAndServe(":8080", nil)
}
