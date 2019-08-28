package internal

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//IndexHandler handles HTTP Requests to the Printern Homepage
func IndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	fmt.Fprintln(w, "Hello World")
}
