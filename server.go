package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Entry struct {
	SessionID uint32
	Alice     User
	Bob       User
}

type Memory struct {
	Mapping map[uint32]*Entry
}

func SendToUser() {
	//TODO
}

var connections = make(map[string]net.Conn)
var publicKeys = make(map[string]rsa.PublicKey)

func StartServer(host string) {

	l, err := net.Listen("tcp", host)
	log.Println("Starting server on:", host)
	// connections := make(map[string]net.Conn)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer l.Close()
	log.Println("Server is running...")

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}

}

// Handles incoming requests.
func handleRequest(conn net.Conn) {

	buff := make([]byte, 1024)
	_, err := conn.Read(buff) // write rec data in buff
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	rec := string(buff)
	log.Println("Got message:", rec)

	if strings.HasPrefix(rec, "REG:") {
		s := strings.Split(rec, ":")
		username := s[1]
		publicKey := rsa.PublicKey{}
		json.Unmarshal([]byte(s[2]), &publicKey)
		log.Printf("Registering user: %v", username)

		connections[username] = conn
		publicKeys[username] = publicKey

		conn.Write([]byte("Server connection successful\n"))
		return
	}

	s := strings.Split(rec, ", ") // message, sender, receiver
	var message, sender, receiver string

	if len(s) == 3 {
		message = s[0]
		sender = s[1]
		receiver = s[2]
	}

	fmt.Println(message, sender, receiver)

	if c := connections[receiver]; c != nil {
		c.Write([]byte(message))
	} else {
		log.Printf("Could not deliver message to user %v because connection was not found\n", receiver)
	}

	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	// conn.Close()
}
