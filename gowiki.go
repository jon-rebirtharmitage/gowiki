package main

import (
	"html/template"
	"net/http"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"regexp"
	"io/ioutil"
	"encoding/json"
)

type Page struct{
	Title string
	Uid int
	Static int
	Nuerons []neuron
}

func loadPage(title string) (*Page, error) {
	moaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "manifold"}
	noaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "paradox"}
	a := mongo_find(moaddr, title)
	if a.Title == "" {
		c := []neuron{}
		f := mongo_find(moaddr, "INDEX")
		mongo_init(moaddr, axion{"INDEX", (f.Uid+1), 0, nil})
		mongo_insertAxion(moaddr, axion{strings.ToLower(title), (f.Uid+1), 0, nil})
		return &Page{title, (f.Uid+1), 0, c}, errors.New("Need to create this page it does not exist.")
	}
	c := []neuron{}
	for i := range a.Synapse{
		c = append(c, mongo_export(noaddr, a.Synapse[i]))
	}
	//c := mongo_multiexport(noaddr, a.Uid)
	return &Page{title, a.Uid, 0, c}, nil
}

func loadParadox(uid string) (*Page, error) {
	noaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "paradox"}
	u, _ := strconv.Atoi(uid)
	a := mongo_locate(noaddr, u)
	fmt.Println(a)
	return &Page{a[0].Title, a[0].Uid, 0, a}, nil
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
	p, _ := loadPage(title)
	renderTemplate(w, "edit", p)
}

func editsmallHandler(w http.ResponseWriter, r *http.Request, uid string) {
	p, _ := loadParadox(uid)
	renderTemplate(w, "editsmall", p)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "editsmall.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|editsmall|view)/([a-zA-Z0-9]+)$")

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

func Save(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow- Origin", "*") 
	moaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "manifold"}
	noaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "paradox"}
	body, _ := ioutil.ReadAll(r.Body)
	var jsonData = []byte(body)
	var t neuron
	json.Unmarshal(jsonData, &t)
	f := mongo_find(moaddr, "LINDEX")
	mongo_init(moaddr, axion{"LINDEX", (f.Uid+1), 0, nil})
	g := mongo_find(moaddr, t.Tags[0])
	g.Synapse = append(g.Synapse, (f.Uid+1))
	mongo_init(moaddr, axion{t.Tags[0], g.Uid, 0, g.Synapse})
	t.Uid = (f.Uid + 1)
	t.Synapse = append(t.Synapse, t.Uid)
	mongo_insert(noaddr, t)
}

func SmallSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow- Origin", "*") 
	noaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "paradox"}
	body, _ := ioutil.ReadAll(r.Body)
	var jsonData = []byte(body)
	var t neuron
	json.Unmarshal(jsonData, &t)
	mongo_update(noaddr, t)
}


func forwardHandler(w http.ResponseWriter, r *http.Request){
	a := strings.Split(r.URL.String(), "/")
	if a[1] == "" {
		http.Redirect(w, r, "/view/" + "gowiki", http.StatusFound)
	}else{
		http.Redirect(w, r, "/view/" + a[1], http.StatusFound)	
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/editsmall/", makeHandler(editsmallHandler))
	http.HandleFunc("/process/", Save)
	http.HandleFunc("/subprocess/", SmallSave)
	http.HandleFunc("/", forwardHandler)
	http.Handle("/css/",http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/js/",http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/img/",http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	http.ListenAndServe(":8085", nil)
}