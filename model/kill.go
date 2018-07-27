package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type KillId struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	EsiID    int           `bson:"esi_id" json:"esi_id"`
	EsiHash  string        `bson:"esi_hash" json:"esi_hash"`
	EsiValue float64       `bson:"esi_value" json:"esi_value"`
	ZkbValue float64       `bson:"zkb_value" json:"zkb_value"`
	MerValue float64       `bson:"mer_value" json:"mer_value"`
	Killmail Killmail      `bson:"killmail" json:"killmail"`
}

type Killmail struct {
	Attackers     []Attacker `bson:"attackers" json:"attackers"`
	KillmailID    int        `bson:"killmail_id" json:"killmail_id"`
	KillmailTime  time.Time  `bson:"killmail_time" json:"killmail_time"`
	SolarSystemID int        `bson:"solar_system_id" json:"solar_system_id"`
	Victim        Victim     `bson:"victim" json:"victim"`
	Zkb           ZkbData    `bson:"zkb" json:"zkb"`
}

type Item struct {
	Flag           int `bson:"flag" json:"flag"`
	ItemTypeID     int `bson:"item_type_id" json:"item_type_id"`
	QuantityDroppd int `bson:"quantity_dropped" json:"quantity_dropped"`
	Singleton      int `bson:"singleton" json:"singleton"`
}

type Position struct {
	X float64 `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
	Z float64 `bson:"z" json:"z"`
}

type Victim struct {
	AllianceID    int      `bson:"alliance_id" json:"alliance_id"`
	CharacterID   int      `bson:"character_id" json:"character_id"`
	CorporationID int      `bson:"corporation_id" json:"corporation_id"`
	DamageTaken   int      `bson:"damage_taken" json:"damage_taken"`
	Items         []Item   `bson:"items" json:"items"`
	Position      Position `bson:"position" json:"position"`
	ShipTypeID    int      `bson:"ship_type_id" json:"ship_type_id"`
}

type ZkbData struct {
	LocationID  int     `bson:"locationID" json:"locationID"`
	Hash        string  `bson:"hash" json:"hash"`
	FittedValue float64 `bson:"fittedValue" json:"fittedValue"`
	TotalValue  float64 `bson:"totalValue" json:"totalValue"`
	Points      int     `bson:"points" json:"points"`
	NPC         bool    `bson:"npc" json:"npc"`
	Solo        bool    `bson:"solo" json:"solo"`
	Awox        bool    `bson:"awox" json:"awox"`
	Esi         string  `bson:"esi" json:"esi"`
	URL         string  `bson:"url" json:"url"`
}

type Attacker struct {
	AllianceID     int     `bson:"alliance_id" json:"alliance_id"`
	CharacterID    int     `bson:"character_id" json:"character_id"`
	CorporationID  int     `bson:"corporation_id" json:"corporation_id"`
	DamageDone     int     `bson:"damage_done" json:"damage_done"`
	FactionID      int     `bson:"faction_id" json:"faction_id"`
	FinalBlow      bool    `bson:"final_blow" json:"final_blow"`
	SecurityStatus float32 `bson:"security_status" json:"security_status"`
	ShipTypeID     int     `bson:"ship_type_id" json:"ship_type_id"`
	WeaponTypeID   int     `bson:"weapon_type_id" json:"weapon_type_id"`
}
