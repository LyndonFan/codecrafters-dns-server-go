package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	calledArgs := os.Args
	if len(calledArgs) != 3 {
		fmt.Println("Usage: --resolver <ip>:<port>")
		return
	}

	address := calledArgs[2]
	resolverAddress, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Printf("Failed to resolve address %s: %v", address, err)
		return
	}

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	resolverConn, err := net.ListenUDP("udp", resolverAddress)
	if err != nil {
		fmt.Println("Failed to bind to resolver address:", err)
		return
	}
	defer resolverConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		receivedPacket := PacketFromBytes(buf[:size])

		responseQuestions := receivedPacket.Questions

		answers := make([]Answer, len(receivedPacket.Questions))
		for i, q := range receivedPacket.Questions {
			intermediatePacket := PacketFromQAs([]Question{q}, []Answer{})
			intermediatePacket.Header.Identifier = uint16(i)
			responsebuf := make([]byte, 512)
			_, err = resolverConn.Write(intermediatePacket.AsBytes())
			if err != nil {
				fmt.Println("Failed to send intermediate request:", err)
				break
			}
			intermediateSize, _, err := resolverConn.ReadFromUDP(responsebuf)
			if err != nil {
				fmt.Println("Failed to receive intermediate response:", err)
				break
			}
			intermediateResponse := PacketFromBytes(responsebuf[:intermediateSize])
			answers[i] = intermediateResponse.Answers[0]
		}

		responsePacket := PacketFromQAs(responseQuestions, answers)
		responsePacket.Header.Identifier = receivedPacket.Header.Identifier
		if receivedPacket.Header.OpCode == 0x00 {
			responsePacket.Header.ResponseCode = 0x00
		} else {
			responsePacket.Header.ResponseCode = 0x04
		}

		response := responsePacket.AsBytes()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
