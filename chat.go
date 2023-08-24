package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	SERVER_HOST = "172.20.10.2"
	SERVER_PORT = "8080"
)

func encrypt(message string, key int) string {
	encrypted := ""
	for _, char := range message {
		encrypted += fmt.Sprintf("%c", (int(char)+key)%256)
	}
	return encrypted
}

func decrypt(message string, key int) string {
	decrypted := ""
	for _, char := range message {
		decrypted += fmt.Sprintf("%c", (int(char)-key+256)%256)
	}
	return decrypted
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Use socket <ip>")
	}

	SERVER_HOST = os.Args[1]
	go Server()
	Client()

}

func Server() {
	server, err := net.Listen("tcp", "172.20.10.2:"+SERVER_PORT)
	if err != nil {
		log.Panic(err)
	}
	defer server.Close()

	fmt.Println("\nlistening on" + "172.20.10.2:" + SERVER_PORT)
	fmt.Println("Esperando pelo cliente...")

	//parado esperando um cliente

	connection, err := server.Accept()
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("client connected")
	//defer connection.Close()
	for {
		buffer := make([]byte, 1024)
		mLen, err := connection.Read(buffer)
		if err == nil {
			encrypted := string(buffer[:mLen])
			decrypted := decrypt(encrypted, 3)
			fmt.Println(decrypted)
		}
		time.Sleep(time.Second * 1)
	}
}

func Client() {
	fmt.Println("connecting on " + SERVER_HOST + ":" + SERVER_PORT)
	connection, err := net.Dial("tcp", SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer connection.Close()
	fmt.Println("Server connected")

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		encrypted := encrypt(text, 3)

		_, err = connection.Write([]byte(encrypted))
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}
	}
}
