package db

import (
	"github.com/arkezar/evehub/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

type Evehub struct {
	Database *mgo.Database
}

func connect(config *viper.Viper) *mgo.Database {
	conn, err := mgo.Dial(config.GetString("database.host"))

	if err != nil {
		panic(err)
	}

	db := conn.DB(config.GetString("database.name"))

	err = db.Login(config.GetString("database.user"), config.GetString("database.pass"))

	if err != nil {
		panic(err)
	}
	return db
}

func New(config *viper.Viper) *Evehub {
	return &Evehub{Database: connect(config)}
}

func (db Evehub) GetUnfetchedKillmails() *mgo.Iter {
	return db.Database.C("kills").Find(bson.M{"killmail.victim.corporation_id": bson.M{"$exists": false}}).Iter()
}

func (db Evehub) UpdateKillmail(km model.KillId) {
	err := db.Database.C("kills").Update(bson.M{"esi_id": km.EsiID}, km)
	if err != nil {
		panic(err)
	}
}
