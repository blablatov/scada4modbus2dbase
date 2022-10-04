package modbus2rtu

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/goburrow/modbus"
)

func TestModbusRtu(t *testing.T) {
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
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadDiscreteInputs(15, 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("ReadDiscreteInputs %#v\r\n", results)
	log.Println("\nResponse of slave RTU modbus server:", results)
}
