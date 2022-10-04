// Modbus main package for parse and send data to MongoDB or MSSQL.
// Основной Modbus-пакет парсер для отправки данных в MongoDB или MSSQL.
package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/blablatov/scada4modbus2dbase/chatbotclient"
	"github.com/blablatov/scada4modbus2dbase/gmodbus2tcp"
	"github.com/blablatov/scada4modbus2dbase/modbus2mgo"
	"github.com/blablatov/scada4modbus2dbase/modbus2rtu"
	"github.com/blablatov/scada4modbus2dbase/modbus2tcp"
)

type DataType struct {
	protocolType string
	methodType   string
	dataType     string
	rdataType    uint16
}

const (
	DsnMongo = "mongodb://localhost:27017/testdb"
)

func main() {
	http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe("localhost:8080", nil))
	log.Fatal(http.ListenAndServeTLS("localhost:8443", "server.crt", "server.key", nil))
}

// this handler is returning component path of URL.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	fmt.Fprintf(w, "Server is listening on 8443. Go to https://127.0.0.1:8443")

	// Formating data of structure. Заполнение структуры. //TODO slice.
	i := strings.Index(r.URL.Path, ":") // get index of symbol ":". Получить индекс первого символа ":" строки
	sd := &DataType{}                   // get types of structure. Получить и присвоить значения типам структуры.
	sd.protocolType = r.URL.Path[:i]    // get slice of string before symbol ":". Получить значение до первого символа ":".
	substr := r.URL.Path[i+1:]          // get slice of string after symbol ":". Получить подстроку substr после первого символа ":".

	t := strings.Index(substr, ":") // get index symbol ":" in substr. Получить индекс 1-го символа ":" в подстроке substr.
	sd.methodType = substr[:t]      // get slice of string before symbol ":". Получить значение до первого символа ":" substr.
	fmt.Println("Protocol type:", strings.TrimPrefix(sd.protocolType, "/"))

	switch sd.protocolType {
	// Switching type of modbus tcp. Переключатель на modbus tcp.
	case "/modbus_tcp":
		// Switching to function of modbus tcp. Переключатель на функцию modbus tcp.
		switch sd.methodType {
		case "WriteSingleCoil":
			////////////////////////////////////////////////////////////////////
			// Sending data on Modbus via method WriteSingleCoil of interface.
			// Отправка данных через вызов метода WriteSingleCoil интерфейса.
			sd.dataType = substr[t+1:] // get slice of string after symbol ":". Получить значение после первого символа ":" substr.
			var wdb bool
			switch sd.dataType {
			case "true":
				wdb = true
			case "false":
				wdb = false
			default:
				fmt.Println("\nData type is incorrect for WriteSingleCoil method")
			}
			// Formating data of structure Modbus. Заполнение структуры.
			md := modbus2tcp.ModbusData{
				SwitchMethodType:    sd.methodType,
				WriteSingleCoilData: wdb,
			}
			start := time.Now()
			// Calling an interface method.
			// Вызов метода WriteSingleCoil интерфейса.
			var d modbus2tcp.Modbuser = md
			reswrite, err := d.WriteSingleCoil()
			if err != nil {
				log.Fatalf("Error of method: %v", err)
			}
			fmt.Println("Result of request via interface method WriteSingleCoil: ", reswrite)
			secs := time.Since(start).Seconds()
			fmt.Printf("%.2fs Request execution time via method of interface\n", secs)

			////////////////////////////////////////////////////////////////////
			// Sending data on Modbus via goroutine.
			// Отправка данных методом ReadCoils через горутину.
			start2 := time.Now()
			// Channel to data exchange for methods of write. Канал обмена данными для методов записи.
			cw := make(chan bool)
			// Synchronization of goroutines. Синхронизация горутин.
			var wgt sync.WaitGroup
			wgt.Add(1) // Counter of goroutines. Значение счетчика горутин.
			go gmodbus2tcp.GWriteSingleCoil(wgt, sd.methodType, wdb, cw)
			// Getting data from goroutine. Получение данных из канала горутины.
			log.Println("\nResult of request via method WriteSingleCoil of goroutine: ", <-cw)
			secs2 := time.Since(start2).Seconds()
			fmt.Printf("%.2fs Request execution time via method of goroutine\n", secs2)
			// Wait of counter. Ожидание счетчика.
			wgt.Wait()
			close(cw)
		case "ReadCoils":
			////////////////////////////////////////////////////////////////////
			// Reading data on Modbus via method ReadCoils of interface.
			// Чтение данных через вызов метода ReadCoils интерфейса.
			sd.dataType = substr[t+1:] // get slice of string after symbol ":". Получить значение после первого символа ":" substr.
			dt, err := strconv.ParseUint(sd.dataType, 10, 64)
			if err != nil {
				fmt.Println(err.Error())
			}
			// Formating data of structure Modbus. Заполнение структуры.
			rd := modbus2tcp.ModbusData{
				SwitchMethodType: sd.methodType,
				ReadCoilsData:    uint16(dt),
			}
			start3 := time.Now()
			// Calling an interface method.
			// Вызов метода ReadCoils интерфейса.
			var d modbus2tcp.Modbuser = rd
			result, err := d.ReadCoils()
			if err != nil {
				log.Fatalf("Error of method: %v", err)
			}
			fmt.Println("Result of request via interface method ReadCoils: ", result)

			strchat := intsToString(result) // Call of func for convert to string. Преобразование в строку.
			cd := chatbotclient.ChatData{
				ModbusData: strchat,
			}
			// Send of scada-data to clients via chatbot. Отправка scada-данных клиентам через чат-бот.
			// Calling an interface ChatClient method. Вызов метода ChatClient интерфейса.
			var c chatbotclient.ChatUser = cd
			c.ChatClient()
			secs3 := time.Since(start3).Seconds()
			fmt.Printf("%.2fs Request execution time via method ReadCoils of interface\n", secs3)

			///////////////////////////////////////
			// Writting data of modbus to MongoDB via method SendMongo of interface.
			// Запись данных Modbus в MongoDB через метод SendMongo интерфейса.
			if len(result) == 0 {
				fmt.Println(err.Error())
			}
			stres := intsToString(result) // Call of func for convert to string. Преобразование в строку.
			// Formating data of structure Modbus. Заполнение структуры.
			dm := modbus2mgo.ModbusMongo{
				SensorType:     "Dallas1",
				SensModbusData: stres,
			}
			start4 := time.Now()
			// Calling an interface method.
			// Вызов метода SendMongo интерфейса.
			var s modbus2mgo.ModbusMonger = dm
			resreq, err := s.SendMongo(DsnMongo)
			if err != nil {
				log.Fatalf("Error of method: %v", err)
			}
			fmt.Println("Result of request via interface method ReadCoils: ", resreq)
			secs4 := time.Since(start4).Seconds()
			fmt.Printf("%.2fs Request execution time via method SendMongo of interface\n", secs4)

			////////////////////////////////////////////////////////////////////
			// Read data method GReadCoils via goroutine.
			// Чтение данных методом GReadCoils через горутину.
			sd.dataType = substr[t+1:] // get slice of string after symbol ":". Получить значение после первого символа ":" substr.
			dg, err := strconv.ParseUint(sd.dataType, 10, 64)
			if err != nil {
				fmt.Println(err.Error())
			}
			gd := uint16(dg)
			start2 := time.Now()
			// Channel to data exchange for methods of read. Канал обмена данными для методов чтения.
			cr := make(chan []byte)
			// Synchronization of goroutines. Синхронизация горутин.
			var wgr sync.WaitGroup
			wgr.Add(1) // Counter of goroutines. Значение счетчика горутин.
			go gmodbus2tcp.GReadCoils(wgr, sd.methodType, gd, cr)
			// Getting data from goroutine. Получение данных из канала горутины.
			log.Println("\nResult of request via method ReadCoils of goroutine: ", <-cr)
			secs2 := time.Since(start2).Seconds()
			fmt.Printf("%.2fs Request execution time via method GReadCoils of goroutine\n", secs2)
			// Wait of counter. Ожидание счетчика.
			wgr.Wait()
			close(cr)
		default:
			fmt.Println("\nFunction not found")
		}
	case "modbus_rtu":
		start3 := time.Now()
		cr := make(chan []byte) // Channel to data exchange. Канал для обмена данными.
		var wgr sync.WaitGroup  // Synchronization of goroutines. Синхронизация горутин.
		wgr.Add(1)              // Counter of goroutines. Значение счетчика горутин.
		go modbus2rtu.ModbusRtu(wgr, cr)
		// Getting data from goroutine. Получение данных из канала горутины.
		log.Println("\nModbus RTU data: ", <-cr)
		secs3 := time.Since(start3).Seconds()
		fmt.Printf("%.2fs Request execution time RTU via goroutine\n", secs3)
		// Wait of counter. Ожидание счетчика.
		go func() {
			wgr.Wait()
		}()
	default:
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "Data for page not found: %s\n", r.URL.Path)
	}
}

// Function like fmt.Sprint(), but it insert commas. Функция аналогична fmt.Sprint(), + добавляет запятые.
func intsToString(values []byte) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	//buf.WriteRune('') // For any rune in UTF-8 encoding. Для произвольной руны в кодировке UTF-8.
	for i, v := range values {
		if i > 0 {
			buf.WriteString(", ") // It separate values with a comma. Разделяет значения строки запятой.
		}
		fmt.Fprintf(&buf, "%d", v)
	}
	buf.WriteByte(']')
	return buf.String()
}
