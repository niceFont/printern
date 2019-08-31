package internal

import (
	"github.com/thedevsaddam/renderer"
)

//Rnderer initiates the Render stuct and exposes it to the rest of the application
var Rnderer *renderer.Render

/*
//Chrome Devtools Protocol
var Chrome context.Context

//ChromeCancel Function
var ChromeCancel context.CancelFunc */

func init() {

	/* Chrome, ChromeCancel = chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	) */

	rndOptions := renderer.Options{
		ParseGlobPattern: "./web/*.html",
	}

	Rnderer = renderer.New(rndOptions)
}
