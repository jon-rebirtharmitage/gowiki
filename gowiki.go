package main

import (
	"html/template"
	"net/http"
	"errors"
	"fmt"
	"strings"
	"regexp"
	//"io"
	"io/ioutil"
	"encoding/json"
	//"log"
)

type Page struct {
	Title string
	Body  string
}

func (p *Page) save() error {
// 	moaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "test"}
// 	a := mongo_export(moaddr, mogvalue)
// 	if len(a) == 0 {
// 		mongo_insert(moaddr, movalue)
// 	} 
// 	mongo_update(moaddr, movalue)
	return nil
}

func loadPage(title string) (*Page, error) {
	moaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "index"}
	t := neuron{0,0, title, "", nil, nil}
	a := mongo_export(moaddr, t)
	if len(a) == 0 {
		return &Page{}, errors.New("Need to create this page it does not exist.")
	}
	return &Page{Title: a[0].Title, Body: a[0].Content}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: body}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func RESTFUL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow- Origin", "*") 
	moaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "test"}
	body, _ := ioutil.ReadAll(r.Body)
	
	var jsonData = []byte(body)
	var t neuron
	json.Unmarshal(jsonData, &t)
	mongo_insert(moaddr, t)
	fmt.Println(t.Title)
	http.Redirect(w, r, "/view/" + t.Title, http.StatusFound)
}

func forwardHandler(w http.ResponseWriter, r *http.Request){
	a := strings.Split(r.URL.String(), "/")
	http.Redirect(w, r, "/view/" + a[1], http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/process/", RESTFUL)
	http.HandleFunc("/", forwardHandler)
	http.Handle("/css/",http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/js/",http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/img/",http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	http.ListenAndServe(":8085", nil)
}