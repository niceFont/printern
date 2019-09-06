package internal

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	internal "github.com/user/printern/internal/app"
)

//IndexHandler handles HTTP Requests to the Printern Homepage
func IndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if pusher, ok := w.(http.Pusher); ok {
		if err := pusher.Push("./web/css/app.css", nil); err != nil {
			log.Printf("Failed to push %v\n", err)
		}
	}
	t, err := template.ParseFiles("./web/html/index.html")

	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, nil)
}

func PrinterHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	pr := &internal.PrinterRequest{}
	bodyB, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error while reading request!", http.StatusBadRequest)
	}

	err = json.Unmarshal(bodyB, pr)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pr.JobTitle)
	res := internal.DispatchCrawlers(pr)

	response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, "Error while Parsing Response", http.StatusInternalServerError)
		return
	}
	log.Printf("Sending Response: %s", response)

	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}
