package server

import (
	"log"
	"syscall"
)

var (
	PORT       = 8080
	ADDR       = [4]byte{127, 0, 0, 1}
	BACKLOG    = 128
	MAXMSGSIZE = 8000
)

func RunServer() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)

	if err != nil {
		log.Fatal("Error creating server socket: ", err)
	}

	err = syscall.Bind(fd, &syscall.SockaddrInet4{Port: PORT, Addr: ADDR})

	if err != nil {
		log.Fatal("Error binding server socket: ", err)
	}

	err = syscall.Listen(fd, BACKLOG)

	if err != nil {
		log.Fatal("Error listening to server socket: ", err)
	}

	newFd, _, err := syscall.Accept(fd)

	if err != nil {
		log.Fatal("Error accepting socket: ", err)
	}

	msg := make([]byte, MAXMSGSIZE)

	numBytes, _, err := syscall.Recvfrom(newFd, msg, 0)

	if err != nil {
		log.Fatal("Error receiving message: ", err)
	}

	log.Printf("Message received (%d bytes): %s\n", numBytes, msg)

	err = syscall.Close(newFd)

	if err != nil {
		log.Fatal("Error closing new socket: ", err)
	}

	err = syscall.Close(fd)

	if err != nil {
		log.Fatal("Error closing original socket: ", err)
	}
}
