package main

import (
	"errors"
	"strconv"

	"github.com/kataras/iris/v12"
	Calc "github.com/samdotme/refi-calc"
)

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html").Reload(true))

	app.HandleDir("/assets", iris.Dir("./assets"))

	app.Get("/", index)
	app.Post("/results", handleForm)

	app.Listen(":8080")
}

func index(ctx iris.Context) {
	data := iris.Map{
		"Title":      "Home Refinance Calculator",
		"FooterText": "Footer contents",
		"Message":    "Main contents",
	}

	ctx.ViewLayout("layouts/main")
	ctx.View("index", data)
}

type refiForm struct {
	MonthlyPayment string `form:"monthly_payment"`
	CurrentAmount  string `form:"current_amount"`
	CurrentRate    string `form:"current_rate"`
	RefiRate       string `form:"refi_rate"`
	ClosingCosts   string `form:"closing_costs"`
}

type refiData struct {
	MonthlyPayment float64
	CurrentAmount  float64
	CurrentRate    float64
	RefiRate       float64
	ClosingCosts   float64
}

func handleForm(ctx iris.Context) {
	var form refiForm

	err := ctx.ReadForm(&form)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	data, err := convertFormToData(form)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	curPaydownPeriod := Calc.CalculatePaydownPeriod(data.CurrentAmount, data.MonthlyPayment, data.CurrentRate)

	refiAmount := data.CurrentAmount + data.ClosingCosts
	refiPaydownPeriod := Calc.CalculatePaydownPeriod(refiAmount, data.MonthlyPayment, data.RefiRate)

	ctx.JSON(iris.Map{
		"inputData":                  data,
		"currentPaydownPeriodMonths": curPaydownPeriod,
		"currentPaydownPeriodYears":  curPaydownPeriod / 12,
		"refiPaydownPeriodMonths":    refiPaydownPeriod,
		"refiPaydownPeriodYears":     refiPaydownPeriod / 12,
		"shouldRefi":                 refiPaydownPeriod < curPaydownPeriod,
	})
}

func convertFormToData(form refiForm) (refiData, error) {
	monthlyPayment, err := strconv.ParseFloat(form.MonthlyPayment, 64)
	if err != nil {
		return refiData{}, errors.New("cannot parse monthly payment")
	}

	curAmt, err := strconv.ParseFloat(form.CurrentAmount, 64)
	if err != nil {
		return refiData{}, errors.New("cannot parse current amount")
	}

	curRate, err := strconv.ParseFloat(form.CurrentRate, 64)
	if err != nil {
		return refiData{}, errors.New("cannot parse current rate")
	}

	refiRate, err := strconv.ParseFloat(form.RefiRate, 64)
	if err != nil {
		return refiData{}, errors.New("cannot parse refinance rate")
	}

	closingCosts, err := strconv.ParseFloat(form.ClosingCosts, 64)
	if err != nil {
		return refiData{}, errors.New("cannot parse closing costs")
	}

	data := refiData{
		MonthlyPayment: monthlyPayment,
		CurrentAmount:  curAmt,
		CurrentRate:    curRate * 0.01,
		RefiRate:       refiRate * 0.01,
		ClosingCosts:   closingCosts,
	}

	return data, nil
}
