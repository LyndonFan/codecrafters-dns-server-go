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

		payloadHeader := HeaderFromBytes(buf[:12])
		fmt.Println(payloadHeader)

		payloadQuestions := make([]Question, payloadHeader.QuestionCount)
		qStartIndex := 12
		for i := 0; i < int(payloadHeader.QuestionCount); i++ {
			payloadQuestions[i], qStartIndex = QuestionFromBytes(buf[qStartIndex:], qStartIndex)
		}

		responseQuestions := make([]Question, len(payloadQuestions))
		for i, q := range payloadQuestions {
			responseQuestions[i] = Question{
				Name:  q.Name,
				Type:  1,
				Class: 1,
			}
		}

		answers := make([]Answer, len(payloadQuestions))
		for i, q := range payloadQuestions {
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
			Identifier:        payloadHeader.Identifier,
			QR:                true,
			OpCode:            payloadHeader.OpCode,
			RecursionDesired:  payloadHeader.RecursionDesired,
			QuestionCount:     uint16(len(responseQuestions)),
			AnswerRecordCount: uint16(len(answers)),
		}
		if payloadHeader.OpCode == 0x00 {
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
