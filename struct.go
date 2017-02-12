package main

import (
	"html/template"
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
  Uid, ContentType int
  Title, Parent string
	Content template.HTML
  Tags []string
  Synapse []int
}

type axion struct{
  Title string
	Uid int
	Static int
  Synapse []int
}