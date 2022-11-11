// Modbus main package for parse and send data to MongoDB or MSSQL.
// Основной Modbus-пакет парсер для отправки данных в MongoDB или MSSQL.
package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sort"
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
	mu           sync.Mutex
	protocolType string
	methodType   string
	dataType     string
	rdataType    uint16
}

// Anonymous field. Composition for secure access to types and methods of the modbus2tcp package.
// Анонимное поле. Композиция для безопасного доступа к типам и методам пакета modbus2tcp.
type embtypes struct {
	modbus2tcp.ModbusData
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

	////////////////////////////////////////////////////////////////////
	// Parsing strings of input for requst to modbus a slave device.
	// Парсинг строк запроса к modbus-slave устройству.
	sd := &DataType{}
	counts := make(map[string]int)
	//Slice for save all keys from mapping. Срез для хранения всех ключей мапы.
	datakeys := make([]string, 0, len(counts))
	for _, line := range strings.Split(string(r.URL.Path), ":") {
		counts[line]++
		log.Println(line) // Checks data of mapping.
	}
	// Sorts keys to list and to sets values in order.
	// Сортировка ключей для перечисления и присваивания значений по порядку.
	for countkeys := range counts {
		datakeys = append(datakeys, countkeys)
	}
	sort.Strings(datakeys)
	for _, countkeys := range datakeys {
		fmt.Printf("\nCountkeys: %v\nCounts: %v\n", countkeys, counts[countkeys])
		if countkeys != "" {
			if sd.protocolType == "" {
				sd.mu.Lock()
				sd.protocolType = countkeys
				sd.mu.Unlock()
			} else {
				if sd.dataType == "" {
					sd.mu.Lock()
					sd.dataType = countkeys
					sd.mu.Unlock()
				} else {
					sd.mu.Lock()
					sd.methodType = countkeys
					sd.mu.Unlock()
				}
			}
		}
	}
	// Output for test. Тестовый вывод данных.
	log.Println("protocolType: ", sd.protocolType)
	log.Println("dataType: ", sd.dataType)
	log.Println("methodType: ", sd.methodType)
	fmt.Println("\nProtocol type check:", strings.TrimPrefix(sd.protocolType, "/"))

	switch sd.protocolType {
	// Switching type of modbus tcp. Переключатель на modbus tcp.
	case "/modbus_tcp":
		// Switching to function of modbus tcp. Переключатель на функцию modbus tcp.
		switch sd.methodType {
		case "WriteSingleCoil":
			////////////////////////////////////////////////////////////////////
			// Sending data on Modbus via method WriteSingleCoil of interface.
			// Отправка данных через вызов метода WriteSingleCoil интерфейса.
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
			reswrite := d.WriteSingleCoil()
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
			dt, err := strconv.ParseUint(sd.dataType, 10, 64)
			if err != nil {
				fmt.Println(err.Error())
			}

			// Option one.
			// Calling an interface method via struct embedding.
			// Вызов метода ReadCoils интерфейса, через встроенную структуру.
			start3 := time.Now()
			var w embtypes
			w.ReadCoilsData = uint16(dt)
			w.SwitchMethodType = sd.methodType
			fmt.Println(w)

			result, err := embtypes.ReadCoils(w)
			if err != nil {
				log.Fatalf("Error of method: %v", err)
			}
			fmt.Println("Result of request via interface method: ", result)

			// Option two.
			// Formating data of structure Modbus. Заполнение структуры.
			/*rd := modbus2tcp.ModbusData{
				SwitchMethodType: sd.methodType,
				ReadCoilsData:    uint16(dt),
			}
			// Calling an interface method.
			// Вызов метода ReadCoils интерфейса.
			var d modbus2tcp.Modbuser = rd
			result, err := d.ReadCoils()
			if err != nil {
				log.Fatalf("Error of method: %v", err)
			}
			fmt.Println("Result of request via interface method ReadCoils: ", result)*/

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
			p := recover()
			if len(result) == 0 {
				log.Println("Panic, internal error, data of answed not got. recover()")
				panic(p)
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
			resreq := s.SendMongo(DsnMongo)
			fmt.Println("Result of request via interface method ReadCoils: ", resreq)
			secs4 := time.Since(start4).Seconds()
			fmt.Printf("%.2fs Request execution time via method SendMongo of interface\n", secs4)

			////////////////////////////////////////////////////////////////////
			// Read data method GReadCoils via goroutine.
			// Чтение данных методом GReadCoils через горутину.
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
