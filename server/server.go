package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
)

var conList ConList

// ConList creates a connection list
type ConList struct {
	lock sync.Mutex
	list []net.Conn
}

// Append appends an item to the list
func (list *ConList) Append(c net.Conn) {
	list.lock.Lock()
	list.list = append(list.list, c)
	list.lock.Unlock()
}

// Clear clears list.
func (list *ConList) Clear() {
	list.lock.Lock()
	list.list = nil
	list.lock.Unlock()
}

// Print prints the content of the list
func (list *ConList) Print() {
	list.lock.Lock()
	fmt.Print(list.list)
	list.lock.Unlock()
}

var playList PlayList

// PlayList creates a list with the plays
type PlayList struct {
	lock sync.Mutex
	list []int64
}

// Append appends an item to the list
func (list *PlayList) Append(v int64) {
	list.lock.Lock()
	list.list = append(list.list, v)
	list.lock.Unlock()
}

// Print prints the content of the list
func (list *PlayList) Print() {
	list.lock.Lock()
	fmt.Print(list.list)
	list.lock.Unlock()
}

// Clear clears list.
func (list *PlayList) Clear() {
	list.lock.Lock()
	list.list = nil
	list.lock.Unlock()
}

func printMove(move int64) string {
	if move == 0 {
		return "Rock"
	} else if move == 1 {
		return "Paper"
	} else {
		return "Scissors"
	}
}

func handlePlays(userPlay int64, player int64) {
	playList.Append(userPlay)
}

func handleConnection(conn net.Conn, player int64) {
	var user int64
	err := binary.Read(conn, binary.LittleEndian, &user)
	if err != nil {
		fmt.Println("Failed to receive move", err)
		os.Exit(3)
	}
	if user > 2 || user < 0 {
		fmt.Println("Invalid Play", user)
		os.Exit(4)
	}
	handlePlays(user, player)

	if player == 1 {
		handleScore(err)
	}
}

func handleScore(err error) {
	move1 := playList.list[0]
	move2 := playList.list[1]

	result1 := "Your move:" + printMove(move1) + "Play of the opponent:" + printMove(move2)
	result2 := "Your move:" + printMove(move2) + "Play of the opponent:" + printMove(move1)

	difference := (move1 - move2) % 3
	if difference < 0 {
		difference += 3
	}

	switch difference {
	case 0:
		result1 += "Draw! \n"
		result2 += "Draw! \n"
	case 1:
		result1 += "You Won! \n"
		result2 += "You Lose! \n"
	case 2:
		result1 += "You Lose! \n"
		result2 += "You Won! \n"
	}

	for i := 0; i < 2; i++ {
		if i == 0 {
			_, err = fmt.Fprint(conList.list[i], result1)
		}
		if i == 1 {
			_, err = fmt.Fprint(conList.list[i], result2)
		}
		if err != nil {
			fmt.Println("Failed to send result", err)
			os.Exit(5)
		}
	}

}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server", err)
		os.Exit(1)
	}

	var playerCount int64
	playerCount = -1
	fmt.Println("Available server")

	for {
		conn, err := ln.Accept()
		conList.Append(conn)
		playerCount = playerCount + 1

		if playerCount > 1 {
			playerCount = 0
			playList.Clear()
			conList.Clear()
		}

		if err != nil {
			fmt.Println("Error accepting connection:", err)
			os.Exit(2)
		}
		handleConnection(conn, playerCount)
	}
}
