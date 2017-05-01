package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Connection failed:", err)
		os.Exit(1)
	}

	var user int64
	str1 := "Rock (0), Paper (1) or Scissors (2)?"

	fmt.Println(str1)
	fmt.Scanf("%d", &user)

	if user > 2 || user < 0 {
		fmt.Println("Invalid move \nThe possible moves are: Stone (0)," +
			"Paper (1) or Scissors (2)")
		conn.Close()
		os.Exit(2)
	}
	fmt.Println(user)

	err = binary.Write(conn, binary.LittleEndian, &user)
	if err != nil {
		fmt.Println("Failed to send move:", err)
		conn.Close()
		os.Exit(3)
	}
	var result string
	result, err = bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Println("Failed to receive result:", err)
		conn.Close()
		os.Exit(4)
	}
	fmt.Print(result)
	conn.Close()
}
