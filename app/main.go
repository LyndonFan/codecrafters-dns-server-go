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

	resolverConn, err := net.DialUDP("udp", nil, resolverAddress)
	if err != nil {
		fmt.Println("Failed to bind to resolver address:", err)
		return
	}
	defer resolverConn.Close()
	fmt.Println("resolverConn.LocalAddr()", resolverConn.LocalAddr())

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

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		fmt.Printf("Received %d bytes from %s\n", size, source)

		packet := PacketFromBytes(buf[:size])
		fmt.Printf("Received packet:\n%v\n", packet)

		receivedQuestions := packet.Questions
		fmt.Printf("Received %d question(s)\n", len(receivedQuestions))

		answers := make([]Answer, len(packet.Questions))
		for i, q := range packet.Questions {
			packet.Questions = []Question{q}
			intermediateResponse, err := sendRequest(resolverConn, &packet)
			if err != nil {
				fmt.Println("Failed to send intermediate request:", err)
				continue
			}
			answers[i] = intermediateResponse.Answers[0]
		}

		packet.Questions = receivedQuestions
		packet.Answers = answers
		if packet.Header.OpCode != 0x00 {
			packet.Header.ResponseCode = 0x04
		}
		packet.Header.AnswerRecordCount = uint16(len(answers))
		packet.Header.QR = true

		response := packet.AsBytes()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func sendRequest(
	resolverConn *net.UDPConn,
	packet *Packet,
) (*Packet, error) {
	bytes := packet.AsBytes()
	nBytes, err := resolverConn.Write(bytes)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return nil, err
	}
	fmt.Printf("Sent packet %d bytes to resolver\n", len(bytes))
	fmt.Printf("Sent %d bytes written\n", nBytes)

	responseBuf := make([]byte, 512)
	responseSize, _, err := resolverConn.ReadFromUDP(responseBuf)
	if err != nil {
		fmt.Println("Failed to receive response:", err)
		return nil, err
	}
	fmt.Printf("Received %d bytes from resolver\n", responseSize)

	responsePacket := PacketFromBytes(responseBuf[:responseSize])
	return &responsePacket, nil
}
