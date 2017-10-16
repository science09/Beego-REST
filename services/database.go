package services

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
)

type DataBase struct {
	s       *mgo.Session
	name    string
	session *mgo.Database
}

func (db *DataBase) Connect() {
	db.s = service.Session()
	session := *db.s.DB(db.name)
	db.session = &session
}

func newDBSession(name string) *DataBase {
	var db = DataBase{name: name}
	db.Connect()
	return &db
}

func checkAndInitServiceConnection() {
	if service.baseSession == nil {
		service.MongoURL = beego.AppConfig.String("DBPath")
		maxPool, err := beego.AppConfig.Int("DBMaxPool")
		if err != nil {
			panic(err)
		}
		err = service.New(maxPool)
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	//init method to start db connection
	beego.Info("init mongo DB...")
	checkAndInitServiceConnection()
}
