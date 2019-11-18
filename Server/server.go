package main

import (
	"bufio"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
)

func check(err error, message string) {
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", message)
}

type ClientJob struct {
	name    string
	numbers string
	conn    net.Conn
}

func isPrime(number int64) bool {

	if big.NewInt(number).ProbablyPrime(0) {
		return true
	}

	return false
}

func countDigits(i int64) (count int) {
	if i == 0 {
		count = 1
		return count
	}

	for i != 0 {

		i /= 10
		count++
	}
	return count
}

func generateResponses(clientJobs chan ClientJob, maxSizeArray int) {
	for {
		clientJob := <-clientJobs

		fmt.Printf("A client with the name %s has connected \n", clientJob.name)

		fmt.Printf("Server recieved a request from %s \n", clientJob.name)

		fmt.Printf("Server is processing the request... \n")

		numbersArray := strings.Fields(clientJob.numbers)

		if len(numbersArray) > maxSizeArray {

			clientJob.conn.Write([]byte("The array size exceeds the maximum allowed\n"))

		} else {
			count := 0

			for _, element := range numbersArray {

				intElement, err := strconv.ParseInt(element, 10, 64)

				if err != nil {
					fmt.Println("Non-numeric value found in array")
					count = -1
					break
				}

				if isPrime(intElement) {
					count += countDigits(intElement)
				}

			}

			if count < 0 {

				clientJob.conn.Write([]byte("Your input is not an array of numbers\n"))

			} else {

				countString := strconv.FormatInt(int64(count), 10)
				clientJob.conn.Write([]byte("The total number of digits from all prime numbers" +
					" in the array is: " + countString + "\n"))

			}

			fmt.Printf("Server has sent a response to %s\n", clientJob.name)
		}

	}
}

func setConf(filePath string) (nbClients int, maxSizeArray int) {

	fd, err := os.Open(filePath)

	if err != nil {
		panic(fmt.Sprintf("open %s: %v", filePath, err))

	}

	_, err = fmt.Fscanf(fd, "%d", &nbClients)

	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Scan Failed %s: %v", filePath, err))

	}

	_, err = fmt.Fscanf(fd, "%d", &maxSizeArray)

	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Scan Failed %s: %v", filePath, err))

	}

	fd.Close()
	return
}

func main() {
	nbClients, sizeArray := setConf("Config/config.txt")
	clientJobs := make(chan ClientJob)
	go generateResponses(clientJobs, sizeArray)

	ln, err := net.Listen("tcp", ":8080")
	check(err, "Server is ready.")

	for {
		if nbClients >= 1 {
			conn, err := ln.Accept()
			check(err, "Accepted connection.")

			if err == nil {

				nbClients--

				go func() {
					buf := bufio.NewReader(conn)

					for {
						name, err := buf.ReadString('\n')

						if err != nil {
							fmt.Printf("Client disconnected.\n")
							nbClients++
							break
						}

						name = strings.Join(strings.Fields(name), "")

						numbers, err := buf.ReadString('\n')

						if err != nil {
							fmt.Printf("Client disconnected.\n")
							nbClients++
							break
						}

						clientJobs <- ClientJob{name, numbers, conn}
					}
				}()
			}

		}

	}

}
