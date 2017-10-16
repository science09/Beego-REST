package services

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
)

type Service struct {
	baseSession *mgo.Session
	queue       chan int
	MongoURL    string
	Open        int
}

var service Service

func (s *Service) New(maxPool int) error {
	var err error
	s.queue = make(chan int, maxPool)
	for i := 0; i < maxPool; i++ {
		s.queue <- 1
	}
	s.Open = 0
	s.baseSession, err = mgo.Dial(s.MongoURL)
	//username := beego.AppConfig.String("DBUser")
	//password := beego.AppConfig.String("DBPassword")
	//err = s.baseSession.DB("admin").Login(username, password)
	//if err != nil {
	//	beego.Error("Can not connect to MongoDb: " + err.Error())
	//	panic(err)
	//}

	return err
}

func (s *Service) Session() *mgo.Session {
	<-s.queue
	s.Open++
	return s.baseSession.Copy()
}

func (s *Service) Close(c *Collection) {
	c.db.s.Close()
	s.queue <- 1
	s.Open--
}

type Collection struct {
	db      *DataBase
	name    string
	Session *mgo.Collection
}

func NewCollectionSession(name string) *Collection {
	db_name := beego.AppConfig.String("DBName")
	var c = Collection{db: newDBSession(db_name), name: name}
	c.Connect()
	return &c
}

func (c *Collection) Connect() {
	session := *c.db.session.C(c.name)
	c.Session = &session
}

func (c *Collection) Close() {
	service.Close(c)
}
