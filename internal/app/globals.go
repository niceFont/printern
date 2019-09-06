package internal

import (
	"github.com/thedevsaddam/renderer"
)

//Rnderer initiates the Render stuct and exposes it to the rest of the application
var Rnderer *renderer.Render

var Keywords []string

func init() {

	Keywords = []string{"php", "javascript", "react", "c", "c++", "java", "c#", "angular", "nodejs"}

	rndOptions := renderer.Options{
		ParseGlobPattern: "./web/*.html",
	}

	Rnderer = renderer.New(rndOptions)
}
