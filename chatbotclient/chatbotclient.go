// Demo chat client for data exchange with a chat server via tls.
// Демо чат-клиент для обмена данными с чат-сервером по tls.
package chatbotclient

import (
	"crypto/tls"
	"io"
	"log"
	"os"
)

type ChatUser interface {
	ChatClient()
}

// Structure of data modbus. Структура данных modbus.
type ChatData struct {
	ModbusData string
}

// Creating chat channel. Создание чат-канала.
func (d ChatData) ChatClient() {
	log.SetPrefix("Client event: ")
	log.SetFlags(log.Lshortfile)

	// Allows not reg a sertificaties, it's for demo mode. Разрешение не регать сертификаты, для демо.
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	//conn, err := net.Dial("tcp", "localhost:8888")
	conn, err := tls.Dial("tcp", "localhost:4443", conf)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors. Примечание: игнорирует ошибки.
		log.Println("done")
		done <- struct{}{} // signal the main goroutine. Сигнал главной go-подпрограмме.
	}()
	// Main programm read type of struct, send data to server. Главная программа считывает тип структуры и отправляет данные на сервер.
	mustCopy(conn, d.ModbusData) // Sends the row data to the server. Отправляет данные строки на сервер.
	// Reads the program's standard input and sends the data to the server. Считывает стандартный ввод программы и отправляет данные на сервер.
	//mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // Wait for background goroutine to finish. Ожидание завершения фоновой go-подпрограммы.
}

// Func for send data via type of structure. Для отправки данных через тип структуры.
func mustCopy(dst io.Writer, src string) {
	if _, err := io.WriteString(dst, src); err != nil {
		log.Fatal(err)
	}
}

// Func for send data via io.Reader. Func для отправки данных через io.Reader.
/*func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}*/
