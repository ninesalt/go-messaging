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

type StandardMessage struct {
	Header string
}

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
		go handleRequest(conn)
	}

}

// Handles incoming requests.
func handleRequest(conn net.Conn) {

	scanner := bufio.NewScanner(conn)
	defer conn.Close()

	for scanner.Scan() {
		content := scanner.Bytes()
		parsed := StandardMessage{}
		json.Unmarshal(content, &parsed)

		if err := scanner.Err(); err != nil {
			log.Printf("error reading connection: %v\n", err)
			break
		}

		switch parsed.Header {
		case "REG":
			handleRegistration(conn, content)
		case "GETPKEY":
			handlePublicKeyRetrieval(conn, content)
		}
	}

}

// handleRegistration is the main handler for new users announcing their
// username and public key to the server
func handleRegistration(conn net.Conn, content []byte) {

	type RegisterMessage struct {
		Username  string
		PublicKey rsa.PublicKey
	}

	r := RegisterMessage{}
	json.Unmarshal(content, &r)
	username := strings.ToLower(r.Username)
	connections[username] = conn
	publicKeys[username] = r.PublicKey
	log.Println("Registered user: ", username)
	conn.Write([]byte("User registration successful\n"))
}

func handlePublicKeyRetrieval(conn net.Conn, content []byte) {

	type GetPublicKeyMessage struct {
		Username string // the partner's username
	}

	message := GetPublicKeyMessage{}
	json.Unmarshal(content, &message)
	u := strings.ToLower(message.Username)
	pkey := publicKeys[u]
	encoder := json.NewEncoder(conn)
	encoder.Encode(pkey)
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
