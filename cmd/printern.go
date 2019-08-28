package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	internal "github.com/user/printern/internal/pkg"
)

func main() {

	router := httprouter.New()
	router.GET("/", internal.IndexHandler)

	http.ListenAndServe(":3000", router)
}
