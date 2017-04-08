package main

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Node struct {
	Id int64

	//url地址
	Url string `xorm:"notnull unique"`

	//標題
	Name string

	Gid string
}

var g_engine *xorm.Engine

func init() {
	engine, e := xorm.NewEngine("sqlite3", "my.db")
	if e != nil {
		log.Fatal(e)
	}
	e = engine.Ping()
	if e != nil {
		log.Fatal(e)
	}
	engine.ShowSQL(false)

	bean := &Node{}
	if ok, e := engine.IsTableExist(bean); e != nil {
		log.Fatal(e)
	} else if ok {
		engine.Sync2(bean)
	} else {
		if e = engine.CreateTables(bean); e != nil {
			log.Fatal(e)
		}
		if e = engine.CreateIndexes(bean); e != nil {
			log.Fatal(e)
		}
		if e = engine.CreateUniques(bean); e != nil {
			log.Fatal(e)
		}
	}

	g_engine = engine
}
func GetEngine() *xorm.Engine {
	return g_engine
}
func NewSession() *xorm.Session {
	return g_engine.NewSession()
}
