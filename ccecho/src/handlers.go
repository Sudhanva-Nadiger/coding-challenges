package src

import (
	"fmt"
	"io"
	"log"
	"net"
)

func HandleTcpConnection(network string, port int) {
	l, err := net.ListenTCP(network, &net.TCPAddr{
		Port: port,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()
	fmt.Printf("TCP server listening on port %d\n", port)

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			defer c.Close()
			defer fmt.Println("Client closed", c.RemoteAddr().String())

			fmt.Printf("Accepted TCP connection from: %s\n", c.RemoteAddr().String())

			io.Copy(c, c)
		}(conn)
	}
}

func HandleUdpConnection(network string, port int) {
	addr := &net.UDPAddr{
		Port: port,
	}

	conn, err := net.ListenUDP(network, addr)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("UDP server listening on port %d\n", port)

	buffer := make([]byte, 1024)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP: %v\n", err)
			continue
		}

		fmt.Printf("Received datagram from: %s (%d bytes)\n", clientAddr.String(), n)

		// Echo the data back to the client
		_, err = conn.WriteToUDP(buffer[:n], clientAddr)
		if err != nil {
			log.Printf("Error writing UDP: %v\n", err)
			continue
		}
	}
}
