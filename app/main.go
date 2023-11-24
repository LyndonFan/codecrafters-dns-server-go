package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	calledArgs := os.Args
	if len(calledArgs) != 3 {
		fmt.Println("Usage: --resolver <ip>:<port>")
		return
	}

	address := calledArgs[2]
	IP_PATTERN := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}$`)
	if !IP_PATTERN.MatchString(address) {
		fmt.Println("Invalid IP address:", address)
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

		responseQuestions := make([]Question, len(receivedPacket.Questions))
		for i, q := range receivedPacket.Questions {
			responseQuestions[i] = Question{
				Name:  q.Name,
				Type:  1,
				Class: 1,
			}
			fmt.Printf("Question %d: %s\n", i, q.Name)
		}

		answers := make([]Answer, len(receivedPacket.Questions))
		for i, q := range receivedPacket.Questions {
			answers[i] = Answer{
				Name:     q.Name,
				Type:     1,
				Class:    1,
				TTL:      60,
				RDLength: 4,
				RDData:   []byte{0x08, 0x08, 0x08, 0x08},
			}
		}

		header := Header{
			Identifier:        receivedPacket.Header.Identifier,
			QR:                true,
			OpCode:            receivedPacket.Header.OpCode,
			RecursionDesired:  receivedPacket.Header.RecursionDesired,
			QuestionCount:     uint16(len(responseQuestions)),
			AnswerRecordCount: uint16(len(answers)),
		}
		if receivedPacket.Header.OpCode == 0x00 {
			header.ResponseCode = 0x00
		} else {
			header.ResponseCode = 0x04
		}

		response := make([]byte, 0, 512)
		response = append(response, header.AsBytes()...)
		for _, question := range responseQuestions {
			response = append(response, question.AsBytes()...)
		}
		for _, answer := range answers {
			response = append(response, answer.AsBytes()...)
		}

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
