package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
)

type Header struct {
	Identifier            uint16
	QR                    bool
	OpCode                uint8
	Authoritative         bool
	Truncation            bool
	RecursionDesired      bool
	RecursionAvailable    bool
	Reserved              uint8
	ResponseCode          uint8
	QuestionCount         uint16
	AnswerRecordCount     uint16
	AuthorityRecordCount  uint16
	AdditionalRecordCount uint16
}

func (h Header) AsBytes() []byte {
	res := make([]byte, 12)
	res[0] = byte(h.Identifier >> 8)
	res[1] = byte(h.Identifier & 0xff)

	if h.QR {
		res[2] |= 1 << 7
	}
	res[2] |= uint8(h.OpCode) << 3
	if h.Authoritative {
		res[2] |= 1 << 2
	}
	if h.Truncation {
		res[2] |= 1 << 1
	}
	if h.RecursionDesired {
		res[2] |= 1
	}

	if h.RecursionAvailable {
		res[3] |= 1 << 7
	}
	res[3] |= h.Reserved << 4
	res[3] |= h.ResponseCode

	res[4] = byte(h.QuestionCount >> 8)
	res[5] = byte(h.QuestionCount & 0xff)

	res[6] = byte(h.AnswerRecordCount >> 8)
	res[7] = byte(h.AnswerRecordCount & 0xff)

	res[8] = byte(h.AuthorityRecordCount >> 8)
	res[9] = byte(h.AuthorityRecordCount & 0xff)

	res[10] = byte(h.AdditionalRecordCount >> 8)
	res[11] = byte(h.AdditionalRecordCount & 0xff)

	return res
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
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

		// Create an empty response
		response := []byte{}
		header := calculateHeaders()
		response = append(response, header...)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func calculateHeaders() []byte {
	header := Header{
		Identifier: 1234,
		QR:         true,
	}

	return header.AsBytes()

	// header := make([]byte, 0, 12)
	// identifier := 1234
	// header = append(header, byte(identifier>>8), byte(identifier&0xff)) // 1234 = 0x04d2

	// qrIndicator := 1
	// oopCode := 0
	// authoritativeAnswer := 0
	// truncation := 0
	// recursionDesired := 0
	// val := (qrIndicator << 7) | (oopCode << 3)
	// val = val | (authoritativeAnswer << 2) | (truncation << 1) | recursionDesired
	// header = append(header, byte(val))

	// recursionAvailable := 0
	// reserved := 0
	// responseCode := 0
	// val = (recursionAvailable << 7) | (reserved << 4) | responseCode
	// header = append(header, byte(val))

	// questionCount := 0
	// header = append(header, byte(questionCount>>8), byte(questionCount&0xff))

	// answerRecordCount := 0
	// header = append(header, byte(answerRecordCount>>8), byte(answerRecordCount&0xff))

	// authorityRecordCount := 0
	// header = append(header, byte(authorityRecordCount>>8), byte(authorityRecordCount&0xff))

	// additionalRecordCount := 0
	// header = append(header, byte(additionalRecordCount>>8), byte(additionalRecordCount&0xff))

	// return header
}
