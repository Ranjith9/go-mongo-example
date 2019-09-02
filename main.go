package main

import (
	"html/template"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"gopkg.in/mgo.v2"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

type Profile struct {
	FirstName string
	LastName  string
}

func main() {
	http.HandleFunc("/", foo)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/users", getUsers)
	http.ListenAndServe(":3000", nil)
}

func foo(w http.ResponseWriter, req *http.Request) {

	f := req.FormValue("first")
	l := req.FormValue("last")

	mine := Profile{f, l}
	if f != "" {
		dbInsert(mine)
	}
	err := tpl.ExecuteTemplate(w, "index.gohtml", Profile{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}

func dbInsert(p Profile) {

	session, err := mgo.Dial("db:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("profiles").C("users")
	err = c.Insert(&Profile{p.FirstName, p.LastName})
	if err != nil {
		log.Fatal(err)
	}

}

func getUsers(w http.ResponseWriter, req *http.Request) {
	
	session, err := mgo.Dial("db:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	var results []Profile

	c := session.DB("profiles").C("users")

	err1 := c.Find(nil).All(&results)
	if err1 != nil {
		panic(err1)
	}

	b, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
  
	json.Unmarshal(b, &results)

	output, _ := json.MarshalIndent(results, "", " ")
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}
