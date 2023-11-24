package main

import (
	"fmt"
	"strings"
)

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

func AnswerFromBytes(data []byte, startIndex int) (Answer, int) {
	a := Answer{}
	i := startIndex
	a.Name, i = GetName(data, i)
	a.Type = uint16(data[i])<<8 | uint16(data[i+1])
	a.Class = uint16(data[i+2])<<8 | uint16(data[i+3])
	a.TTL = uint32(data[i+4])<<24 | uint32(data[i+5])<<16 | uint32(data[i+6])<<8 | uint32(data[i+7])
	a.RDLength = uint16(data[i+8])<<8 | uint16(data[i+9])
	a.RDData = make([]byte, a.RDLength)
	copy(a.RDData, data[i+10:i+10+int(a.RDLength)])
	return a, i + 10 + int(a.RDLength)
}
