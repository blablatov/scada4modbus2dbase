// Demo ChatBot server for any clients chat.
// Демо ChatBot сервер для любого клиентского чата.
package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

type client chan<- string // an outgoing message channel. Канал исходящих сообщений

var (
	// Clobal channels. Глобальные каналы.
	entering = make(chan client) // Connect of clients. Подключение клиентов.
	leaving  = make(chan client) // Disconnect of clients. Отключение клиентов.
	messages = make(chan string) // All incoming clients messages. Все входящие сообщения клиентов.
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients. Все подключенные клиенты
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all clients' outgoing message channels.
			// Широковещательное входящее сообщение во все каналы исходящих сообщений для клиентов.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

// Creates a new client outbound channel and do broadcast via entering channel.
// Создает новый исходящий канал клиента и делает трансляцию через канал entering.
func handleConn(conn net.Conn) {
	ch := make(chan string)   // outgoing client messages. Исходящие сообщения клиентов.
	go clientWriter(conn, ch) // Handler-broadcaster for any client. Обработчик широковещания для клиентов.

	who := conn.RemoteAddr().String()
	ch <- "Device: " + who
	messages <- who + " in touch"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err(). Примечание: игнорируем потенциальные ошибки input.Err().
	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

// Func handler-broadcaster. Функция широковещателя обработчика.
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors. // Примечание: игнорируем ошибки сети.
	}
}

func main() {
	log.SetPrefix("Server event: ")
	log.SetFlags(log.Lshortfile)

	cert, err := tls.LoadX509KeyPair("chatbotserver.crt", "chatbotserver.key")
	if err != nil {
		log.Println(err)
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	//listener, err := net.Listen("tcp", "localhost:8888")
	listener, err := tls.Listen("tcp", "localhost:4443", config)
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
