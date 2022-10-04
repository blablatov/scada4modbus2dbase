package chatbotclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"testing"
)

func TestChatData(t *testing.T) {
	var tests = []struct {
		ModbusData string
	}{
		{"0"},
		{"1,"},
		{"123"},
		{"1999"},
		{"2001"},
	}

	var prevModbusData string
	for _, test := range tests {
		if test.ModbusData != prevModbusData {
			fmt.Printf("\n%s\n", test.ModbusData)
			prevModbusData = test.ModbusData
		}
	}
}

func TestChatClient(t *testing.T) {
	// Allows not reg a sertificaties, it's for demo mode. Разрешение не регать сертификаты, для демо.
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", "localhost:4443", conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connect is successful", conn)

	var dst io.Writer
	dst = new(bytes.Buffer) // Like interface io.Writer. Соответствует интерфейсу io.Writer.
	var src = "123"         // Test data for write to standart stream out. Для теста записи в стандартный поток вывода.
	log.Println("Test data for write to standard output: ", src)
	if _, err := io.WriteString(dst, src); err != nil {
		log.Fatal(err)
	}
	log.Println("done")
	log.Println("Test written data to standard output: ", dst)
	conn.Close()
}

// Func for send data via type of structure. Для отправки данных через тип структуры.
func BenchmarkChatClient(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 3; i++ {
		var dst io.Writer
		dst = new(bytes.Buffer) // Like interface io.Writer. Соответствует интерфейсу io.Writer.
		var src = "123"         // Test data for write to standart stream out. Для теста записи в стандартный поток вывода.
		log.Println("Test data for write to standard output: ", src)
		if _, err := io.WriteString(dst, src); err != nil {
			log.Fatal(err)
		}
		log.Println("done")
		log.Println("Test written data to standard output: ", dst)
	}
}
