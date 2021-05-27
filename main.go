package main

import "github.com/kataras/iris/v12"

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
	CurrentAmount string `form:"current_amount"`
	CurrentRate   string `form:"current_rate"`
	CurrentTerm   string `form:"current_term"`
	RefiRate      string `form:"refi_rate"`
	RefiTerm      string `form:"refi_term"`
	ClosingCosts  string `form:"closing_costs"`
}

func handleForm(ctx iris.Context) {
	var refiData refiForm

	err := ctx.ReadForm(&refiData)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	ctx.JSON(iris.Map{"Refi": refiData})
}
