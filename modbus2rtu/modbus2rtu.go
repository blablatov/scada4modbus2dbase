// Modbus RTU client-master module.
// Base of idea to https://github.com/goburrow/modbus.git
// Thanks him!
package modbus2rtu

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/goburrow/modbus"
)

func ModbusRtu(wgr sync.WaitGroup, cr chan []byte) {
	defer wgr.Done()
	// Modbus RTU/ASCII master
	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler.BaudRate = 115200
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second

	err := handler.Connect()
	if err != nil {
		fmt.Println(err.Error())
		//panic(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadDiscreteInputs(15, 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("ReadDiscreteInputs %#v\r\n", results)
	log.Println("\nResponse of slave RTU modbus server:", results)
	cr <- results
}
