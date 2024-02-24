package client

import (
	"log"
	"syscall"
)

var (
	PORT = 8080
	ADDR = [4]byte{127, 0, 0, 1}
)

func RunClient() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)

	if err != nil {
		log.Fatal("Error creating client socket: ", err)
	}

	serverAddr := &syscall.SockaddrInet4{Port: PORT, Addr: ADDR}

	err = syscall.Connect(fd, serverAddr)

	if err != nil {
		log.Fatal("Error connecting to server: ", err)
	}

	msg := "Hello from the client!"

	err = syscall.Sendmsg(fd, []byte(msg), nil, serverAddr, syscall.MSG_DONTWAIT)

	if err != nil {
		log.Fatal("Error sending message to server: ", err)
	}

	err = syscall.Close(fd)

	if err != nil {
		log.Fatal("Error closing original socket: ", err)
	}
}
