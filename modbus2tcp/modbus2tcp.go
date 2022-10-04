// Modbus TCP client-master module.
// Base of idea to https://github.com/goburrow/modbus.git
// Thanks dude!
package modbus2tcp

import (
	"fmt"
	"sync"
	"time"

	modbus "github.com/thinkgos/gomodbus/v2"
)

type Modbuser interface {
	WriteSingleCoil() bool
	ReadCoils() ([]byte, error)
}

// Structure of data modbus. Структура данных modbus.
type ModbusData struct {
	ReadCoilsData              uint16
	ReadDiscreteInputsData     uint16
	ReadHoldingRegistersData   uint16
	WriteSingleCoilData        bool
	WriteSingleRegisterData    uint16
	WriteMultipleCoilsData     []byte
	WriteMultipleRegistersData []uint16
	SwitchMethodType           string
}

func (d ModbusData) ReadCoils() ([]byte, error) {
	var (
		mu     sync.Mutex // Protect of var of result. Защита от гонки переменной result
		result []byte
	)
	// Modbus TCP master
	p := modbus.NewTCPClientProvider("localhost:5020", modbus.WithEnableLogger())
	client := modbus.NewClient(p)
	err := client.Connect()
	if err != nil {
		fmt.Println("Connect failed, ", err)
		return nil, err
	}
	defer client.Close()

	fmt.Println("\nStarting of interface")

	switch d.SwitchMethodType {
	case "ReadCoils":
		// Reading the state of the device's digital outputs.
		// Чтение состояния цифровых выходов устройства.
		mu.Lock()
		defer mu.Unlock()
		result, err := client.ReadCoils(1, 0, d.ReadCoilsData) // Code of function - Read Coils: 01 (0x01)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("ReadCoils % x\n", result)
			return result, nil
		}
		time.Sleep(time.Second * 2)

	case "ReadDiscreteInputs":
		// Reading the status of the device's digital inputs.
		// Чтение состояния цифровых входов устройства.
		mu.Lock()
		defer mu.Unlock()
		result, err := client.ReadDiscreteInputs(2, 2, d.ReadCoilsData) // Code of function - Read Discrete Inputs:02 (0х02)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("ReadDiscreteInputs % x\n", result)
			return result, nil
		}
		time.Sleep(time.Second * 2)

	case "ReadHoldingRegisters":
		// Reading the device's general purpose registers.
		// Чтение регистров общего назначения устройства.
		mu.Lock()
		defer mu.Unlock()
		result, err := client.ReadHoldingRegisters(3, 0, d.ReadCoilsData) // Code of function - Read Holding Registers:03 (0х03)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("ReadHoldingRegisters % x\n", result)
		}
		time.Sleep(time.Second * 2)
	default:
		fmt.Println("\nFunction read not found")
	}
	return result, nil
}

func (d ModbusData) WriteSingleCoil() bool {
	// Modbus TCP master
	p := modbus.NewTCPClientProvider("localhost:5020", modbus.WithEnableLogger())
	client := modbus.NewClient(p)
	err := client.Connect()
	if err != nil {
		fmt.Println("Connect failed, ", err)
		return false
	}

	defer client.Close()
	fmt.Println("\nStarting of interface")
	switch d.SwitchMethodType {
	case "WriteSingleCoil":
		// Changing the state single of the discrete outputs (coil) of the device.
		// Изменение состояния одного из дискретных выходов устройства.
		fmt.Println("\nWrite Single Coil")
		err := client.WriteSingleCoil(5, 0, true) // Code of function - Write Single Coil:05 (0х05)
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(time.Second * 2)

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
		fmt.Println("\nFunction write not found")
	}
	return true
}
