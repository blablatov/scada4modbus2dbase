package modbus2mgo

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ModbusMonger interface {
	SendMongo(DsnMongo string) bool
}

type ModbusMongo struct {
	SensorType     string
	SensModbusData string
}

func (md ModbusMongo) SendMongo(DsnMongo string) bool {
	session, err := mgo.Dial(DsnMongo)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// is check name in dBase
	c := session.DB("scadadb").C("sensors")
	chk := ModbusMongo{}
	err = c.Find(bson.M{"sensortype": md.SensorType, "datasensor": md.SensModbusData}).One(&chk)
	if err == nil {
		log.Println("\nName already is to DB, method of interface", err)
		return false
	}
	if err != nil {
		log.Print("\nErr data for write to MongoDB, method of interface: ", err)
		err = c.Insert(&ModbusMongo{md.SensorType, md.SensModbusData})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Sensor was written via method of interface:", md.SensorType,
			"\nData of sensor was written via method of interface:", md.SensModbusData)
	}
	return true
}
