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

func mongo_insert(moaddress MOAddr, movalue MOValue) bool{
  session, err := mgo.Dial(moaddress.session)
  if err != nil {
    return false
  }
  defer session.Close()
  c := session.DB(moaddress.table).C(moaddress.doc)
  c.Insert(bson.M{"title": movalue.Title, "body": movalue.Body})
  return true
}

func mongo_export(moaddress MOAddr, mogvalue MOGValue) []MOValue {
	session, err := mgo.Dial(moaddress.session)
	if err != nil {
			panic(err)
	}
	defer session.Close()
	c := session.DB(moaddress.table).C(moaddress.doc)
	result := []MOValue{}
	iter := c.Find(bson.M{"title": mogvalue.Title}).Iter()
	err = iter.All(&result)
	if err != nil {
			log.Fatal(err)
	}
	return result
}

func mongo_update(moaddress MOAddr, movalue MOValue) bool{
  session, err := mgo.Dial(moaddress.session)
  if err != nil {
    return false
  }
  defer session.Close()
  c := session.DB(moaddress.table).C(moaddress.doc)
  c.Update(bson.M{"title": movalue.Title}, bson.M{"title": movalue.Title, "body": movalue.Body})
  return true
}