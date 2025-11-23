package main

import (
	"ccecho/src"
	"log"
	"os"
)

func Listener(network string, port int) {
	if network == "tcp" {
		src.HandleTcpConnection(network, port)
	} else {
		src.HandleUdpConnection(network, port)
	}

}

func main() {
	args := os.Args[1:]

	network := "tcp"
	port := 7

	if len(args) >= 1 {
		// Remove leading dash if present (e.g., -tcp or -udp)
		protocol := args[0]
		if len(protocol) > 0 && protocol[0] == '-' {
			protocol = protocol[1:]
		}

		if protocol != "tcp" && protocol != "udp" {
			log.Fatal("Invalid protocol: ", protocol)
		}

		network = protocol
	}

	Listener(network, port)
}
