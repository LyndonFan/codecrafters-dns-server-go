package main

import (
	"fmt"
	"net"
	"os"
	"sync"
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

	resolverConn, err := net.DialUDP("udp", nil, resolverAddress)
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
		fmt.Printf("Received %d questions\n", len(responseQuestions))

		answers := make([]Answer, len(receivedPacket.Questions))
		responseChannel := make(chan RequestResult, len(receivedPacket.Questions))
		wg := &sync.WaitGroup{}
		for i, q := range receivedPacket.Questions {
			wg.Add(1)
			intermediatePacket := PacketFromQAs([]Question{q}, []Answer{})
			go sendRequest(i, resolverConn, intermediatePacket.AsBytes(), responseChannel, wg)
		}
		wg.Wait()

		for i := 0; i < len(receivedPacket.Questions); i++ {
			requestResult := <-responseChannel
			answers[requestResult.Index] = requestResult.Answer
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

type RequestResult struct {
	Index  int
	Answer Answer
}

func sendRequest(
	index int,
	resolverConn *net.UDPConn,
	buf []byte,
	responseChannel chan (RequestResult),
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	_, err := resolverConn.Write(buf)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}

	responseBuf := make([]byte, 512)
	responseSize, _, err := resolverConn.ReadFromUDP(responseBuf)
	if err != nil {
		fmt.Println("Failed to receive response:", err)
		return
	}

	responsePacket := PacketFromBytes(responseBuf[:responseSize])
	responseChannel <- RequestResult{index, responsePacket.Answers[0]}
}
