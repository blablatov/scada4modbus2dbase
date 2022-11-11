package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/blablatov/scada4modbus2dbase/chatbotclient"
	"github.com/blablatov/scada4modbus2dbase/gmodbus2tcp"
	"github.com/blablatov/scada4modbus2dbase/modbus2tcp"
)

func TestStrings(t *testing.T) {
	var strtest = []struct {
		protocolType    string
		writeMethodType string
		writeDataType   string
		rdataType       uint16
	}{
		{"modbus_tcp", "ReadCoils", "test", 10},
		{"modbus_rtu", "WriteSingleCoil", ",,,:", 16},
		{"\t", "one\ttwo\tthree\\\\", "one/ttwo/tthree/", 8},
		{"Data for test", "#&U*(()))_+_11234", ">>>>????hgjk", 1},
		{"Yes, no", "567, 008, 10", "or, 16, null", 12345},
	}

	var prevprotocolType string
	for _, test := range strtest {
		if test.protocolType != prevprotocolType {
			fmt.Printf("\n%s\n", test.protocolType)
			prevprotocolType = test.protocolType
		}
	}

	var prevwriteMethodType string
	for _, test := range strtest {
		if test.writeMethodType != prevwriteMethodType {
			fmt.Printf("\n%s\n", test.writeMethodType)
			prevwriteMethodType = test.writeMethodType
		}
	}

	var prevwriteDataType string
	for _, test := range strtest {
		if test.writeDataType != prevwriteDataType {
			fmt.Printf("\n%s\n", test.writeDataType)
			prevwriteDataType = test.writeDataType
		}
	}

	var prevrdataType uint16
	for _, test := range strtest {
		if test.rdataType != prevrdataType {
			fmt.Printf("\n%d\n", test.rdataType)
			prevrdataType = test.rdataType
		}
	}
}

func TestSplit(t *testing.T) {
	s, sep := "a:b:c", ":"
	durl := strings.Split(s, sep)
	if got, want := len(durl), 3; got != want {
		t.Errorf("Split (%q%q) возвращает %d слов, а требуется %d", s, sep, got, want)
	}
}

func BenchmarkReadInterface(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 2; i++ {
		strtest := "10"
		dt, err := strconv.ParseUint(strtest, 10, 64)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("DataType:", dt)
		writeMethodType := "ReadCoils"
		rd := modbus2tcp.ModbusData{
			SwitchMethodType: writeMethodType,
			ReadCoilsData:    uint16(dt),
		}
		var d modbus2tcp.Modbuser = rd
		result, err := d.ReadCoils()
		if err != nil {
			log.Fatalf("Error of method: %v", err)
		}
		fmt.Println("Result of request via interface method ReadCoils: ", result)
	}
}

func BenchmarkWriteInterface(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 2; i++ {
		methodType := "WriteSingleCoil"
		md := modbus2tcp.ModbusData{
			SwitchMethodType:    methodType,
			WriteSingleCoilData: true,
		}
		var d modbus2tcp.Modbuser = md
		resfunc := d.WriteSingleCoil()
		fmt.Println("Result of request via interface method WriteSingleCoil: ", resfunc)
	}
}

func BenchmarkRWGoroutines(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 2; i++ {
		methodType := "ReadCoils"
		dataType := "10"
		dg, err := strconv.ParseUint(dataType, 10, 64)
		if err != nil {
			fmt.Println(err.Error())
		}
		gd := uint16(dg)
		wmethodType := "WriteSingleCoil"
		var wdb bool
		wdb = true
		// Channel to data exchange for methods of read. Канал обмена данными для методов чтения.
		wm := make(chan bool)
		cr := make(chan []byte)
		// Synchronization of goroutines. Синхронизация горутин.
		var wgr sync.WaitGroup
		wgr.Add(1) // Counter of goroutines. Значение счетчика горутин.
		go gmodbus2tcp.GReadCoils(wgr, methodType, gd, cr)
		// Getting data from goroutine. Получение данных из канала горутины.
		log.Println("\nResult of request via method ReadCoils of goroutine: ", <-cr)
		go func() {
			wgr.Wait()
			close(cr)
		}()
		wgr.Add(1)
		go gmodbus2tcp.GWriteSingleCoil(wgr, wmethodType, wdb, wm)
		// Getting data from goroutine. Получение данных из канала горутины.
		log.Println("\nResult of request via method WriteSingleCoil of goroutine: ", <-wm)
		// Wait of counter. Ожидание счетчика.
		go func() {
			wgr.Wait()
			close(wm)
		}()
	}
}

func BenchmarkChatClient(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 2; i++ {
		strchat := "0012" // Call of func for convert to string. Преобразование в строку.
		cd := chatbotclient.ChatData{
			ModbusData: strchat,
		}
		// Send of scada-data to clients via chatbot. Отправка scada-данных клиентам через чат-бот.
		// Calling an interface ChatClient method. Вызов метода ChatClient интерфейса.
		var c chatbotclient.ChatUser = cd
		c.ChatClient()
	}
}
