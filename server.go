package main

import (
	"bufio"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

var connections = make(map[string]net.Conn)
var publicKeys = make(map[string]rsa.PublicKey)

func StartServer(host string) {

	l, err := net.Listen("tcp", host)
	log.Println("Starting server on:", host)

	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()
	log.Println("Server is running...")

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Error accepting: ", err.Error())
		}
		// Handle connections in a new goroutine.
		fmt.Println("got a new connection")
		go handleRequest(conn)
	}

}

// Handles incoming requests.
func handleRequest(conn net.Conn) {

	reader := bufio.NewReader(conn)

	for {
		rec, err := reader.ReadString('\n')

		if err != nil {
			// log.Println("Error reading:", err.Error())
			// conn.Close()
			break
		}

		log.Println("Got message")

		if strings.HasPrefix(rec, "REG:") {
			handleRegistration(conn, rec)
		}

		if strings.HasPrefix(rec, "MSG:") {
			handleMessages(rec)
			return
		}

		if strings.HasPrefix(rec, "GETPKEY:") {
			fmt.Println("i am here")
			handlePublicKeyRetrieval(conn, rec)
			return
		}

		// Send a response back to person contacting us.
		conn.Write([]byte("Message received."))
	}

}

// handleRegistration is the main handler for new users announcing their
// username and public key to the server
func handleRegistration(conn net.Conn, rec string) {
	s := strings.Split(rec, ":")
	username := s[1]
	username = strings.ToLower(username)
	publicKey := &rsa.PublicKey{}
	json.Unmarshal([]byte(s[2]), publicKey)
	connections[username] = conn
	publicKeys[username] = *publicKey
	log.Println("Registered user: ", username)
	conn.Write([]byte("Connection successful\n"))
}

func handlePublicKeyRetrieval(conn net.Conn, rec string) {
	split := strings.Split(rec, ":")
	username := split[1] // the username whose public key is being queried
	fmt.Println("username is", username)
	fmt.Println(publicKeys)
	pkey := publicKeys[username]
	fmt.Println("the key is ", pkey)
	// if pkey != nil {
	p, _ := json.Marshal(pkey)
	conn.Write([]byte(p))
	// }

}

func handleMessages(msg string) error {

	split := strings.Split(msg, ":")
	target := split[1]
	c := connections[target]

	if c == nil {
		log.Printf("Could not find connection object for target user%v\n", target)
		return errors.New("No connection found for target")
	}
	c.Write([]byte(msg))
	return nil
}

func main() {
	host := flag.String("host", "localhost", "Host where the server should listen on")
	port := flag.String("port", "5000", "Port to listen on")
	flag.Parse()

	addr := fmt.Sprintf("%v:%v", *host, *port)
	StartServer(addr)
}
