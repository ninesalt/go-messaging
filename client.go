package main

import (
	"bufio"
	"crypto/rsa"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

// Commands sent to the server
const (
	REGISTER     = "REG"
	GETPUBLICKEY = "GETPKEY"
)

type User struct {
	PubKey        rsa.PublicKey
	PartnerPubKey rsa.PublicKey
	Username      string
}

type RegisterMessage struct {
	Header    string
	Username  string
	PublicKey rsa.PublicKey
}

func CreateUser() User {
	_, pubkey := GenerateKeyPair(2048)
	username := fmt.Sprintf("alice%v", rand.Intn(500))
	log.Println("Your random username is:", username)
	return User{PubKey: *pubkey, Username: username}
}

func (u *User) ConnectToServer(host string, partner string) {
	conn, _ := net.Dial("tcp", host)
	connbuff := bufio.NewReader(conn)
	defer conn.Close()

	// register the user
	u.Register(conn, connbuff)

	// get partner public key if passed as a flag
	if partner != "" {
		u.GetPartnerKey(partner, conn, connbuff)
	}

	if isValidPublicKey(u.PartnerPubKey) {
		log.Println("Ready to chat")
		for {

		}
	}

}

func (u *User) Register(conn net.Conn, connbuff *bufio.Reader) {

	enc := json.NewEncoder(conn)
	registermessage := RegisterMessage{Header: REGISTER,
		Username: u.Username, PublicKey: u.PubKey}

	enc.Encode(registermessage)
	message, _ := connbuff.ReadString('\n') // listen for reply
	log.Print(message)
}

func (u *User) GetPartnerKey(target string, conn net.Conn, connbuff *bufio.Reader) {

	type GetPublicKeyMessage struct {
		Header   string
		Username string // the partner's username
	}

	// send the request to get the public key message
	enc := json.NewEncoder(conn)
	g := GetPublicKeyMessage{Header: GETPUBLICKEY, Username: target}
	log.Printf("Retrieving %v's public key from server", target)
	enc.Encode(g)

	// decode the response
	dec := json.NewDecoder(conn)
	publicKey := rsa.PublicKey{}
	dec.Decode(&publicKey)

	if isValidPublicKey(publicKey) {
		u.PartnerPubKey = publicKey
		log.Println("Partner public key retrieved and saved")
	} else {
		log.Println("Could not retrieve partner's public key")
	}
}

// SendMessage sends a message to a target given the message and the target's username
func SendMessage(msg string, target string) {

}

func isValidPublicKey(pkey rsa.PublicKey) bool {
	empty := rsa.PublicKey{}
	return pkey != empty
}

func main() {
	rand.Seed(time.Now().UnixNano())
	host := "localhost:5000"
	u1 := CreateUser()
	target := flag.String("partner", "", "The username of the user you would like to chat with")
	flag.Parse()
	u1.ConnectToServer(host, *target)
}
