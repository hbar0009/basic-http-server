package server

import (
	"bytes"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	PORT       = 8080
	ADDR       = [4]byte{127, 0, 0, 1}
	BACKLOG    = 128
	MAXMSGSIZE = 8000
	BASEPATH   = "/routes"
)

func RunServer() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	serverFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)

	if err != nil {
		log.Fatal("Error creating server socket: ", err)
	}

	go func() {
		for range sigChan {
			log.Print("SIGINT received. Exiting.")
			syscall.Close(serverFd)
			os.Exit(0)
		}
	}()

	err = syscall.Bind(serverFd, &syscall.SockaddrInet4{Port: PORT, Addr: ADDR})

	if err != nil {
		log.Fatal("Error binding server socket: ", err)
	}

	err = syscall.Listen(serverFd, BACKLOG)

	if err != nil {
		log.Fatal("Error listening to server socket: ", err)
	}

	// fds := syscall.FdSet{Bits: [16]int64{}}

	// syscall.Select()

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

	// todo: determine what kind of request it was, and how to handle it
	// maybe into like a Request struct? or just headers?

	// someHtml := handleGETRequest(msg)
	payload := handleGETRequest(msg)

	response := addHttpHeaders(payload)

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

func addHttpHeaders(payload []byte) []byte {
	var out bytes.Buffer

	out.WriteString("HTTP/1.1 200 OK\n")
	out.WriteString("Content-Type: text/html; charset=utf-8\n")
	out.WriteString("\n")

	out.Write(payload)

	return out.Bytes()
}

func handleGETRequest(msg []byte) []byte {
	buf := bytes.NewBuffer(msg)
	firstLine, err := buf.ReadString('\n')

	if err != nil {
		log.Fatal("Error reading message: ", err)
	}

	//todo: parse path properly
	path := strings.Split(firstLine, " ")[1]
	fullPath := resolvePath(path)

	fileContents, err := os.ReadFile(fullPath)

	if err != nil {
		log.Fatal("Error reading requested file: ", err)
	}

	return fileContents
}

func resolvePath(path string) string {
	//todo: handle dynamic routes

	var sb strings.Builder

	wd, err := os.Getwd()

	if err != nil {
		log.Fatal("Error getting current directory: ", err)
	}

	sb.WriteString(wd)
	sb.WriteString(BASEPATH)
	sb.WriteString(path)

	if path[len(path)-1] == '/' {
		sb.WriteString("index.html")
	}

	return sb.String()
}
