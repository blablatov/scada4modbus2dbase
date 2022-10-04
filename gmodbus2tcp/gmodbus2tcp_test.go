package gmodbus2tcp

import (
	"fmt"
	"testing"

	modbus "github.com/thinkgos/gomodbus/v2"
)

func TestReadCoils(t *testing.T) {
	var gtests = []struct {
		methodType string
		gd         uint16
		wdb        bool
	}{
		{"ReadCoils", 011, true},
		{",#&U*(()))_+_11234", 1000, false},
	}

	var prevmethodType string
	for _, test := range gtests {
		if test.methodType != prevmethodType {
			fmt.Printf("\n%s\n", test.methodType)
			prevmethodType = test.methodType
		}
	}

	var prevgd uint16
	for _, test := range gtests {
		if test.gd != prevgd {
			fmt.Printf("\n%d\n", test.gd)
			prevgd = test.gd
		}
	}

	var prevwdb bool
	for _, test := range gtests {
		if test.wdb != prevwdb {
			fmt.Printf("\n%t\n", test.wdb)
			prevwdb = test.wdb
		}
	}
}

func TestTCPGReadCoils(t *testing.T) {
	var gtests = []struct {
		methodType string
		gd         uint16
	}{
		{"ReadCoils", 011},
		{",#&U*(()))_+_11234", 1000},
	}
	for _, grtest := range gtests {
		// Modbus TCP master
		p := modbus.NewTCPClientProvider("localhost:5020", modbus.WithEnableLogger())
		client := modbus.NewClient(p)
		err := client.Connect()
		if err != nil {
			fmt.Println("Connect failed, ", err)
			return
		}
		defer client.Close()

		fmt.Println("\nStarting of goroutine")
		switch grtest.methodType {
		case "ReadCoils":
			// Reading the state of the device's digital outputs.
			// Чтение состояния цифровых выходов устройства.
			resrc, err := client.ReadCoils(1, 0, grtest.gd) // Code of function - Read Coils: 01 (0x01)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("ReadCoils % x\n", resrc)
			}
		default:
			fmt.Println("\nFunction not found")
		}
	}
}

func TestTCPGWriteSingleCoil(t *testing.T) {
	var gtests = []struct {
		methodType string
		wdb        bool
	}{
		{"WriteSingleCoil", true},
		{",#&U*(()))_+_11234", false},
	}
	for _, gwtest := range gtests {
		// Modbus TCP master
		p := modbus.NewTCPClientProvider("localhost:5020", modbus.WithEnableLogger())
		client := modbus.NewClient(p)
		err := client.Connect()
		if err != nil {
			fmt.Println("Connect failed, ", err)
		}
		defer client.Close()

		fmt.Println("\nStarting of interface")
		switch gwtest.methodType {
		case "WriteSingleCoil":
			// Changing the state single of the discrete outputs (coil) of the device.
			// Изменение состояния одного из дискретных выходов устройства.
			fmt.Println("\nWrite Single Coil")
			err := client.WriteSingleCoil(5, 0, gwtest.wdb) // Code of function - Write Single Coil:05 (0х05)
			if err != nil {
				fmt.Println(err.Error())
			}
		default:
			fmt.Println("\nFunction write not found")
		}
	}
}
