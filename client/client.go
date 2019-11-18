package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Type your client's name: ")

	name, _ := reader.ReadString('\n')
	name = strings.Join(strings.Fields(name), "")

	conn, err := net.Dial("tcp", "127.0.0.1:8080")

	if err != nil {
		fmt.Printf("Cannot connect to server. Try again latter.")
		return
	}

	for {

		fmt.Print("Text to send: ")
		numbers, _ := reader.ReadString('\n')

		_, err := fmt.Fprintf(conn, name+"\n"+numbers)

		if err != nil {
			fmt.Printf("Cannot connect to server. Try again latter.")
			return
		}

		fmt.Printf("%s sent a request with the following content: %s \n ", name, numbers)

		message, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			fmt.Printf("Cannot connect to server. Try again latter.")
			return
		}

		fmt.Print("Message from server: " + message)
	}
}
