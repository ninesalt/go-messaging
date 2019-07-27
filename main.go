package main

func main() {
	host := "localhost:5000"
	// StartServer(host)

	u1 := CreateUser()
	u1.ConnectToServer(host)
}
