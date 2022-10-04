package modbus2mgo

import (
	"fmt"
	"log"
	"testing"

	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
)

const (
	DsnMongo = "mongodb://localhost:27017/testdb"
)

func TestSensData(t *testing.T) {
	var tests = []struct {
		SensorType     string
		SensModbusData []byte
	}{
		{"dallas", []byte("0012")},
		{"dallas,", []byte("23456789")},
		{"dallas_1", []byte("12345000987654")},
		{"#&U*(()))_+_11234", []byte("1234561111111111100000000000")},
		{"Yes, dallas,dallas,", []byte("9000090900000000000000000001000000000000000000000000000000000000002")},
	}

	var prevSensorType string
	for _, test := range tests {
		if test.SensorType != prevSensorType {
			fmt.Printf("\n%s\n", test.SensorType)
			prevSensorType = test.SensorType
		}
	}

	var prevSensModbusData []byte
	for _, test := range tests {
		if test.SensModbusData != nil {
			fmt.Printf("\n%s\n", test.SensModbusData)
			prevSensModbusData = test.SensModbusData
		}
	}

	strSensModbusData := string(prevSensModbusData)

	session, err := mgo.Dial(DsnMongo)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// is check name in dBase
	c := session.DB("scadadb").C("sensors")
	chk := ModbusMongo{}
	err = c.Find(bson.M{"sensortype": prevSensorType, "datasensor": strSensModbusData}).One(&chk)
	if err == nil {
		log.Println("\nName already is to DB, method of interface", err)
	}
	if err != nil {
		log.Print("\nErr data for write to MongoDB, method of interface: ", err)
		err = c.Insert(&ModbusMongo{prevSensorType, strSensModbusData})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Sensor was written via method of interface:", prevSensorType,
			"\nData of sensor was written via method of interface:", strSensModbusData)
	}
}
