package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	internal "github.com/user/printern/internal/pkg"
)

func main() {

	router := httprouter.New()
	router.GET("/", internal.IndexHandler)
	router.POST("/printer", internal.PrinterHandler)
	router.NotFound = http.StripPrefix("/static/", http.FileServer(http.Dir("web")))
	http.ListenAndServe(":3000", router)
}
