// +build appengine

package vktopgroups

import (

	"io"
	"log"
	//"fmt"

	"net/http"
	"text/template"


	//"github.com/martini-contrib/oauth2"
	//"github.com/martini-contrib/sessions"
    "github.com/go-martini/martini"
	//"github.com/martini-contrib/cors"

	
	"github.com/crhym3/go-endpoints/endpoints"

	"core"

	"appengine"

)



func serve404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "Not Found")
}



type Params struct {
	ClientId  string
	RootUrl string

}



func oauth2callback(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	c.Infof("handle oauth2callback")
}




func handleMainPage(w http.ResponseWriter, r *http.Request) {

	params := Params{
		ClientId: core.ClientId,
		RootUrl: core.RootUrl,
	}
		

	var index = template.Must(template.ParseFiles("polymer/index.html"))

	 err := index.Execute(w, params)
     if err != nil {
     	log.Fatalf("template execution: %s", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
     }

}





func init() {


	m := martini.Classic()


	m.Get("/oauth2callback", oauth2callback)

	m.Get("/", handleMainPage)

	
	http.Handle("/", m)
   
	if _, err := core.RegisterService(); err != nil {
		panic(err.Error())
	}

	endpoints.HandleHttp()

	


}



