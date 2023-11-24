package main

import (
	"fmt"
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

func GetNameBytes(data []byte, startIndex int) (string, int) {
	if data[startIndex] == 0 {
		return "", startIndex
	}
	nameBytes := make([]byte, 0, len(data)-4)
	i := startIndex
	for i < len(data)-4 {
		if data[i] == 0x00 {
			i += 1
			break
		}
		if data[i]&0xC0 == 0xC0 {
			pointerIndex := int(data[i]&0x3F)<<8 | int(data[i+1])
			suffixName, _ := GetNameBytes(data, pointerIndex)
			nameBytes = append(nameBytes, []byte(suffixName)...)
			nameBytes = append(nameBytes, byte('.'))
			i += 2
			break
		}
		length := int(data[i])
		nameBytes = append(nameBytes, data[i+1:i+1+length]...)
		i += 1 + length
		nameBytes = append(nameBytes, byte('.'))
	}
	nameBytes = nameBytes[:len(nameBytes)-1]
	return string(nameBytes), i
}

func QuestionFromBytes(data []byte, startIndex int) (Question, int) {
	q := Question{}
	i := startIndex
	q.Name, i = GetNameBytes(data, i)
	q.Type = uint16(data[i])<<8 | uint16(data[i+1])
	q.Class = uint16(data[i+2])<<8 | uint16(data[i+3])
	return q, i + 4
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

func (a *Answer) FromBytes(data []byte) {

}
