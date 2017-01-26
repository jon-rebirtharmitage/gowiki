package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
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

func mongo_insert(moaddress MOAddr, t neuron) bool{
  session, err := mgo.Dial(moaddress.session)
  if err != nil {
    return false
  }
  defer session.Close()
  c := session.DB(moaddress.table).C(moaddress.doc)
  c.Insert(t)
	if err != nil{
		return false
	}
	return true
}

func mongo_export(moaddress MOAddr, t neuron) []neuron {
	session, err := mgo.Dial(moaddress.session)
	if err != nil {
			panic(err)
	}
	defer session.Close()
	c := session.DB(moaddress.table).C(moaddress.doc)
	result := []neuron{}
	iter := c.Find(bson.M{"title": t.Title}).Iter()
	err = iter.All(&result)
	if err != nil {
			log.Fatal(err)
	}
	return result
}

func mongo_update(moaddress MOAddr, t neuron) bool{
  session, err := mgo.Dial(moaddress.session)
  if err != nil {
    return false
  }
  defer session.Close()
  c := session.DB(moaddress.table).C(moaddress.doc)
  c.Update(bson.M{"title": t.Title}, t)
  return true
}