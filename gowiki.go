package main

import (
	"html/template"
	"net/http"
//	"errors"
	"fmt"
	"strconv"
	"strings"
	"regexp"
	"io/ioutil"
	"encoding/json"
)

type Page struct{
	Title string
	Ctitle string
	Uid int
	Static int
	Nuerons []neuron
}

func loadPage(title string) (*Page, error) {
	moaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "pantheon"}
	noaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "hades"}
	a := mongo_find(moaddr, title)
	c := []neuron{}
	for i := range a.Synapse{
		c = append(c, mongo_export(noaddr, a.Synapse[i]))
	}
	fmt.Println(a.Ctitle)
	return &Page{title, a.Ctitle, a.Uid, 0, c}, nil
}

func loadParadox(uid string) (*Page, error) {
	noaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "hades"}
	u, _ := strconv.Atoi(uid)
	a := mongo_locate(noaddr, u)
	fmt.Println(a)
	return &Page{"", a[0].Title, a[0].Uid, 0, a}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
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
	//Allows the call to the RESTFUL API to come from across domains
	w.Header().Set("Access-Control-Allow- Origin", "*") 
	//Configures addr information 
	moaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "pantheon"}
	noaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "hades"}
	body, _ := ioutil.ReadAll(r.Body)
	var jsonData = []byte(body)
	var t neuron
	json.Unmarshal(jsonData, &t)
	f := mongo_find(moaddr, "LINDEX")
	mongo_init(moaddr, axion{"LINDEX", "", (f.Uid+1), 0, nil})
	g := mongo_find(moaddr, t.Tags[0])
	g.Synapse = append(g.Synapse, (f.Uid+1))
	mongo_init(moaddr, axion{t.Tags[0], t.Ctitle,  g.Uid, 0, g.Synapse})
	t.Uid = (f.Uid + 1)
	t.Synapse = append(t.Synapse, t.Uid)
	mongo_insert(noaddr, t)
}

func SmallSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow- Origin", "*") 
	noaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "hades"}
	body, _ := ioutil.ReadAll(r.Body)
	var jsonData = []byte(body)
	var t neuron
	json.Unmarshal(jsonData, &t)
	mongo_update(noaddr, t)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow- Origin", "*") 
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	var jsonData = []byte(body)
	var t search
	json.Unmarshal(jsonData, &t)
	moaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "pantheon"}
	noaddr := MOAddr{"vps.rebirtharmitage.com:21701", "gowiki", "hades"}
	t.Searchterms = strings.ToLower(t.Searchterms)
	t.Searchables = strings.Split(t.Searchterms, " ")
	g := []axion{}
	h := []axion{}
	for i := range t.Searchables{
		h = mongo_seekfind(moaddr, string(t.Searchables[i]))
		for j := range h {
			g = append(g, h[j])
		}
	}
	fmt.Println(t.Searchterms)
	if len(g) == 0 {
		f := mongo_find(moaddr, "INDEX")
		j := CreateSessionID()
		mongo_init(moaddr, axion{"INDEX", "", (f.Uid+1), 0, nil})
		k := axion{j, strings.ToLower(t.Searchterms), (f.Uid+1), 0, nil}
		mongo_insertAxion(moaddr, k)
		js, _ := json.Marshal(k)
		w.Write(js)
	}else{
		c := []neuron{}
		for i := range g[0].Synapse{
			c = append(c, mongo_export(noaddr, g[0].Synapse[i]))
		}
		js, _ := json.Marshal(g[0])
		w.Write(js)
	}
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
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/process/", Save)
	http.HandleFunc("/subprocess/", SmallSave)
	http.HandleFunc("/", forwardHandler)
	http.Handle("/css/",http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/js/",http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/img/",http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	http.ListenAndServe(":8085", nil)
}