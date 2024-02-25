package server

import (
	"bytes"
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
	serverFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)

	if err != nil {
		log.Fatal("Error creating server socket: ", err)
	}

	err = syscall.Bind(serverFd, &syscall.SockaddrInet4{Port: PORT, Addr: ADDR})

	if err != nil {
		log.Fatal("Error binding server socket: ", err)
	}

	err = syscall.Listen(serverFd, BACKLOG)

	if err != nil {
		log.Fatal("Error listening to server socket: ", err)
	}

	clientFd, clientAddr, err := syscall.Accept(serverFd)

	if err != nil {
		log.Fatal("Error accepting socket: ", err)
	}

	msg := make([]byte, MAXMSGSIZE)

	numBytes, _, err := syscall.Recvfrom(clientFd, msg, 0)

	if err != nil {
		log.Fatal("Error receiving message: ", err)
	}

	log.Printf("Message received (%d bytes):\n%s\n", numBytes, msg)

	response := createHtmlResponse()

	err = syscall.Sendmsg(clientFd, response, nil, clientAddr, syscall.MSG_DONTWAIT)

	if err != nil {
		log.Fatal("Error sending message to client: ", err)
	}

	err = syscall.Close(clientFd)

	if err != nil {
		log.Fatal("Error closing client socket: ", err)
	}

	err = syscall.Close(serverFd)

	if err != nil {
		log.Fatal("Error closing server socket: ", err)
	}
}

func createHtmlResponse() []byte {
	var out bytes.Buffer

	out.WriteString("HTTP/1.1 200 OK\n")
	out.WriteString("Content-Type: text/html; chatset=utf-8\n")
	out.WriteString("\n")

	someHtml := "<!doctype html><html><p>Hello world!</p><html>"
	out.WriteString(someHtml)

	return out.Bytes()
}
