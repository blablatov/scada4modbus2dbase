package main

import (
	"crypto/tls"
	"log"
	"reflect"
	"testing"
)

func TestMain(t *testing.T) {
	cert, err := tls.LoadX509KeyPair("chatbotserver.crt", "chatbotserver.key")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Test of certificates: ", cert)
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	//listener, err := net.Listen("tcp", "localhost:8888")
	listener, err := tls.Listen("tcp", "localhost:4443", config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Test of run socket: ", listener)

	err = listener.Close()
	if err != nil {
		log.Print(err)
	}
	log.Println("Test of connect close")
}

func Testbroadcaster(t *testing.T) {
	// all connected clients. Все подключенные клиенты
	clients := map[string]bool{
		"sensor_1": true,
		"sensor_2": true,
		"sensor_3": true,
	}
	type client chan<- string
	tests := []struct {
		hello    string
		bye      string
		live     string
		entering <-chan client
		leaving  <-chan client
		messages <-chan string
	}{
		// TODO: Add test cases.
	}
	for _, tb := range tests {
		t.Run(tb.hello, func(t *testing.T) {
			if got := clients; !reflect.DeepEqual(got, tb.entering) {
				t.Errorf("clients() = %v, entering %v", got, tb.entering)
			}
		})
	}

	for _, tb := range tests {
		t.Run(tb.bye, func(t *testing.T) {
			if got := clients; !reflect.DeepEqual(got, tb.bye) {
				t.Errorf("clients() = %v, entering %v", got, tb.bye)
			}
		})
	}

	for _, tb := range tests {
		t.Run(tb.live, func(t *testing.T) {
			if got := clients; !reflect.DeepEqual(got, tb.live) {
				t.Errorf("clients() = %v, entering %v", got, tb.live)
			}
		})
	}
}
