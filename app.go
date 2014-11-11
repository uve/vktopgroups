// +build appengine

package vktopgroups

import (

	"io"
	"log"

	"net/http"
	"text/template"

	"appengine"

	"github.com/crhym3/go-endpoints/endpoints"

	"controller"

)



func serve404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "Not Found")
}



type Params struct {
	ClientId  string
	RootUrl string
	IsDevAppServer bool

}



func handleMainPage(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "text/html")

	params := Params{
		ClientId: controller.ClientId,
		RootUrl: controller.RootUrl,
		IsDevAppServer : appengine.IsDevAppServer(),
	}
		

	var index = template.Must(template.ParseFiles("view/index.html"))

	 err := index.Execute(w, params)
     if err != nil {
     	log.Fatalf("template execution: %s", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
     }

}



func init() {

	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/fetch/contacts", controller.ContactsFetch)
   
	if _, err := controller.RegisterService(); err != nil {
		panic(err.Error())
	}

	endpoints.HandleHttp()

}



