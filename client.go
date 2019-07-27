package main

import (
	"bufio"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
)

type User struct {
	PubKey   *rsa.PublicKey
	Username string
}

func CreateUser() User {
	_, pubkey := GenerateKeyPair(2048)
	username := fmt.Sprintf("alice%v", rand.Intn(500))
	return User{pubkey, username}
}

func (u *User) ConnectToServer(host string) {
	conn, _ := net.Dial("tcp", host)

	for {
		pkey, _ := json.Marshal(u.PubKey)
		text := fmt.Sprintf("REG:%v:%v", u.Username, pkey)
		fmt.Fprintf(conn, text+"\n")
		message, _ := bufio.NewReader(conn).ReadString('\n') // listen for reply
		fmt.Print("Message from server: " + message)
	}
}

func SendMessage(msg string, target string) {

}
