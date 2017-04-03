package main

import (
	"html/template"
	"math/rand"
	"time"
)

type MOAddr struct {
  session, table, doc string
}

type MOValue struct {
  Title, Body string
}

type MOGValue struct {
	Title string
}

type neuron struct{
  Uid int								`json:"uid"` 
	ContentType int				`json:"contenttype"` 
  Title string					`json:"title"` 
	Ctitle string					`json:"ctitle"` 
	Parent string					`json:"parent"` 
	Content template.HTML `json:"content"` 
  Tags []string					`json:"tags"` 
  Synapse []int					`json:"synapse"` 
}

type axion struct{
  Title string			`json:"title"` 
	Ctitle string			`json:"ctitle"`
	Uid int						`json:"uid"`
	Static int				`json:"static"`
  Synapse []int			`json:"synapse"`
}	

type search struct{
	Searchterms string		`json:"searchterms"`
	Searchables []string	`json:"searchables"`
}

type view struct{
	Viewit string
}

func CreateSessionID() (string){
	source := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 24; i++{
		s = s + string(source[rand.Intn(50)])
	}
	return s
}
