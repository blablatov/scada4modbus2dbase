// Modbus TCP client-master module.
// Base of idea to https://github.com/goburrow/modbus.git
// Thanks dude!
package gmodbus2tcp

import (
	"fmt"
	"sync"
	"time"

	modbus "github.com/thinkgos/gomodbus/v2"
)

func GReadCoils(wgr sync.WaitGroup, methodType string, gd uint16, cr chan []byte) {
	defer wgr.Done()
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
	switch methodType {
	case "ReadCoils":
		// Reading the state of the device's digital outputs.
		// Чтение состояния цифровых выходов устройства.
		resrc, err := client.ReadCoils(1, 0, gd) // Code of function - Read Coils: 01 (0x01)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("ReadCoils % x\n", resrc)
			cr <- resrc
		}
		time.Sleep(time.Second * 2)

	case "ReadDiscreteInputs":
		// Reading the status of the device's digital inputs.
		// Чтение состояния цифровых входов устройства.
		resrdi, err := client.ReadDiscreteInputs(2, 2, gd) // Code of function - Read Discrete Inputs:02 (0х02)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("ReadDiscreteInputs % x\n", resrdi)
			cr <- resrdi
		}
		time.Sleep(time.Second * 2)

	case "ReadHoldingRegisters":
		// Reading the device's general purpose registers.
		// Чтение регистров общего назначения устройства.
		resrhr, err := client.ReadHoldingRegisters(3, 0, 10) // Code of function - Read Holding Registers:03 (0х03)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("ReadHoldingRegisters % x\n", resrhr)
			//chr <- resrhr
		}
		time.Sleep(time.Second * 2)
	default:
		fmt.Println("\nFunction not found")
	}
}

func GWriteSingleCoil(wgt sync.WaitGroup, methodType string, wdb bool, cw chan bool) {
	defer wgt.Done()
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
	switch methodType {
	case "WriteSingleCoil":
		// Changing the state single of the discrete outputs (coil) of the device.
		// Изменение состояния одного из дискретных выходов устройства.
		fmt.Println("\nWrite Single Coil")
		err := client.WriteSingleCoil(5, 0, wdb) // Code of function - Write Single Coil:05 (0х05)
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(time.Second * 2)
		cw <- true

	case "WriteSingleRegister":
		// Writing a value to one of the device's holding registers.
		// Запись значения в один из регистров хранения устройства.
		fmt.Println("\nWrite Single Register")
		err = client.WriteSingleRegister(6, 0, 10) // Code of function - Write Single Register:06 (0х06)
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(time.Second * 2)

	case "WriteMultipleCoils":
		// Changes in the state of several discrete outputs of the device.
		// Изменение состояния нескольких дискретных выходов устройства.
		fmt.Println("\nWrite Multiple Coils")
		err = client.WriteMultipleCoils(15, 0, 8, []byte{0}) // Code of function - Write Multiple Coils:15 (0х0F)
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(time.Second * 2)

	case "WriteMultipleRegisters":
		// Writing values to several (from 1 to 123) sequentially located general-purpose registers (registers of holding).
		// Запись значений в несколько (от 1 до 123) последовательно расположенных регистров общего назначения (регистров хранения).
		fmt.Println("\nWrite Multiple Registers")
		err = client.WriteMultipleRegisters(16, 0, 8, []uint16{0, 3, 0, 2, 0, 2, 3, 0}) // Code of function - Write Multiple Registers: 16 (0х10)
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(time.Second * 2)
	default:
		fmt.Println("\nFunction not found")
		cw <- false
	}
}
