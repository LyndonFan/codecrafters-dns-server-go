package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

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

		questions := []Question{}
		questions = append(questions, Question{
			Name:  "codecrafters.io",
			Type:  1,
			Class: 1,
		})

		answers := []Answer{}
		answers = append(answers, Answer{
			Name:     "codecrafters.io",
			Type:     1,
			Class:    1,
			TTL:      60,
			RDLength: 4,
			RDData:   []byte{0x08, 0x08, 0x08, 0x08},
		})

		payloadHeader := Header{}
		payloadHeader.FromBytes(buf[:12])

		header := Header{
			Identifier:        payloadHeader.Identifier,
			QR:                true,
			OpCode:            payloadHeader.OpCode,
			RecursionDesired:  payloadHeader.RecursionDesired,
			QuestionCount:     uint16(len(questions)),
			AnswerRecordCount: uint16(len(answers)),
		}
		if payloadHeader.OpCode == 0x00 {
			header.ResponseCode = 0x00
		} else {
			header.ResponseCode = 0x04
		}

		response := make([]byte, 0, 512)
		response = append(response, header.AsBytes()...)
		for _, question := range questions {
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
