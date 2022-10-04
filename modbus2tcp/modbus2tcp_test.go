package modbus2tcp

import (
	"fmt"
	"sync"
	"testing"

	modbus "github.com/thinkgos/gomodbus/v2"
)

func TestReadCoils(t *testing.T) {
	var tests = []struct {
		ReadCoilsData              uint16
		ReadDiscreteInputsData     uint16
		ReadHoldingRegistersData   uint16
		WriteSingleCoilData        bool
		WriteSingleRegisterData    uint16
		WriteMultipleCoilsData     []byte
		WriteMultipleRegistersData []uint16
		SwitchMethodType           string
	}{
		{012, 2345, 9001, true, 6543, []byte("0012"), []uint16{22}, "any method"},
		{1012, 25, 1100, false, 1000, []byte("110010011"), []uint16{132}, ",#&U*(()))_+_11234"},
	}

	var prevReadCoilsData uint16
	for _, test := range tests {
		if test.ReadCoilsData != prevReadCoilsData {
			fmt.Printf("\n%d\n", test.ReadCoilsData)
			prevReadCoilsData = test.ReadCoilsData
		}
	}

	var prevReadDiscreteInputsData uint16
	for _, test := range tests {
		if test.ReadDiscreteInputsData != prevReadDiscreteInputsData {
			fmt.Printf("\n%d\n", test.ReadDiscreteInputsData)
			prevReadDiscreteInputsData = test.ReadDiscreteInputsData
		}
	}

	var prevReadHoldingRegistersData uint16
	for _, test := range tests {
		if test.ReadHoldingRegistersData != prevReadHoldingRegistersData {
			fmt.Printf("\n%d\n", test.ReadHoldingRegistersData)
			prevReadHoldingRegistersData = test.ReadHoldingRegistersData
		}
	}

	var prevWriteSingleCoilData bool
	for _, test := range tests {
		if test.WriteSingleCoilData != prevWriteSingleCoilData {
			fmt.Printf("\n%t\n", test.WriteSingleCoilData)
			prevWriteSingleCoilData = test.WriteSingleCoilData
		}
	}

	var prevWriteSingleRegisterData uint16
	for _, test := range tests {
		if test.WriteSingleRegisterData != prevWriteSingleRegisterData {
			fmt.Printf("\n%d\n", test.WriteSingleRegisterData)
			prevWriteSingleRegisterData = test.WriteSingleRegisterData
		}
	}

	for _, test := range tests {
		if test.WriteMultipleCoilsData != nil {
			fmt.Printf("\n%s\n", test.WriteMultipleCoilsData)
		}
	}

	for _, test := range tests {
		if test.WriteMultipleRegistersData != nil {
			fmt.Printf("\n%d\n", test.WriteMultipleRegistersData)
		}
	}
}

func TestTCPReadCoils(t *testing.T) {
	var rtests = []struct {
		ReadCoilsData    uint16
		SwitchMethodType string
	}{
		{012, "ReadCoils"},
		{1012, ",#&U*(()))_+_11234"},
	}
	var (
		mu sync.Mutex // Protect of var of result. Защита от гонки переменной result
	)
	for _, rtest := range rtests {
		// Modbus TCP master
		p := modbus.NewTCPClientProvider("localhost:5020", modbus.WithEnableLogger())
		client := modbus.NewClient(p)
		err := client.Connect()
		if err != nil {
			fmt.Println("Connect failed, ", err)
		}
		defer client.Close()

		fmt.Println("\nStarting of interface")
		switch rtest.SwitchMethodType {
		case "ReadCoils":
			// Reading the state of the device's digital outputs.
			// Чтение состояния цифровых выходов устройства.
			mu.Lock()
			defer mu.Unlock()
			result, err := client.ReadCoils(1, 0, rtest.ReadCoilsData) // Code of function - Read Coils: 01 (0x01)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("ReadCoils % x\n", result)
			}
		default:
			fmt.Println("\nFunction read not found")
		}
	}
}

func TestTCPWriteSingleCoil(t *testing.T) {
	var wtests = []struct {
		WriteSingleCoilData bool
		SwitchMethodType    string
	}{
		{true, "WriteSingleCoil"},
		{false, ",#&U*(()))_+_11234"},
	}
	for _, wtest := range wtests {
		// Modbus TCP master
		p := modbus.NewTCPClientProvider("localhost:5020", modbus.WithEnableLogger())
		client := modbus.NewClient(p)
		err := client.Connect()
		if err != nil {
			fmt.Println("Connect failed, ", err)
		}
		defer client.Close()

		fmt.Println("\nStarting of interface")
		switch wtest.SwitchMethodType {
		case "WriteSingleCoil":
			// Changing the state single of the discrete outputs (coil) of the device.
			// Изменение состояния одного из дискретных выходов устройства.
			fmt.Println("\nWrite Single Coil")
			err := client.WriteSingleCoil(5, 0, wtest.WriteSingleCoilData) // Code of function - Write Single Coil:05 (0х05)
			if err != nil {
				fmt.Println(err.Error())
			}
		default:
			fmt.Println("\nFunction write not found")
		}
	}
}
