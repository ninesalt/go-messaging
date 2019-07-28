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

	// for {

	// }
}

func (u *User) Register(conn net.Conn, connbuff *bufio.Reader) {
	pkey, _ := json.Marshal(u.PubKey)
	fmt.Println(pkey)
	text := fmt.Sprintf("%v:%v:%v\n", REGISTER, u.Username, pkey)
	fmt.Fprintf(conn, text)
	message, _ := connbuff.ReadString('\n') // listen for reply
	fmt.Print("Server: " + message)
}

func (u *User) GetPartnerKey(target string, conn net.Conn, connbuff *bufio.Reader) {
	log.Printf("Retrieving %v's public key from server", target)
	text := fmt.Sprintf("%v:%v\n", GETPUBLICKEY, target)
	fmt.Fprintf(conn, text)
	response, _ := connbuff.ReadString('\n')
	fmt.Println("public key:", response)
}

// SendMessage sends a message to a target given the message and the target's username
func SendMessage(msg string, target string) {

}

func main() {
	rand.Seed(time.Now().UnixNano())
	host := "localhost:5000"
	u1 := CreateUser()
	target := flag.String("partner", "", "The username of the user you would like to chat with")
	flag.Parse()
	u1.ConnectToServer(host, *target)
}
