package main

import (
	"fmt"
	"strings"

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

		header := Header{
			Identifier:        1234,
			QR:                true,
			QuestionCount:     uint16(len(questions)),
			AnswerRecordCount: uint16(len(answers)),
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
