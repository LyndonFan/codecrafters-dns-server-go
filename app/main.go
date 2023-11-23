package main

import (
	"fmt"
	"net"
	"strings"
)

type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

func (q Question) AsBytes() []byte {
	res := make([]byte, 0, 4)
	labels := strings.Split(q.Name, ".")
	for _, label := range labels {
		res = append(res, byte(len(label)))
		res = append(res, []byte(label)...)
	}
	res = append(res, byte(0))
	res = append(res, byte(q.Type>>8))
	res = append(res, byte(q.Type&0xff))
	res = append(res, byte(q.Class>>8))
	res = append(res, byte(q.Class&0xff))
	return res
}

type Answer struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RDData   []byte
}

func (a Answer) AsBytes() []byte {
	if len(a.RDData) != int(a.RDLength) {
		fmt.Printf("RDData length (%d) does not match RDLength (%d)\n", len(a.RDData), a.RDLength)
	}
	res := make([]byte, 0, 10)
	labels := strings.Split(a.Name, ".")
	for _, label := range labels {
		res = append(res, byte(len(label)))
		res = append(res, []byte(label)...)
	}
	res = append(res, byte(0))
	res = append(res, byte(a.Type>>8))
	res = append(res, byte(a.Type&0xff))
	res = append(res, byte(a.Class>>8))
	res = append(res, byte(a.Class&0xff))
	for i := 24; i >= 0; i -= 8 {
		res = append(res, byte((a.TTL>>i)&0xff))
	}
	res = append(res, byte(a.RDLength>>8))
	res = append(res, byte(a.RDLength&0xff))
	res = append(res, a.RDData...)
	return res
}

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
