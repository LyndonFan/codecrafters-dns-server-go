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

func GetName(data []byte, startIndex int) (string, int) {
	if data[startIndex] == 0 {
		return "", startIndex + 1
	}
	nameBytes := make([]byte, 0, len(data)-4)
	i := startIndex
	for i <= len(data)-4 {
		if data[i] >= 0xC0 {
			pointerIndex := int(data[i]&0x3F)<<8 | int(data[i+1])
			suffixName, _ := GetName(data, pointerIndex)
			nameBytes = append(nameBytes, []byte(suffixName)...)
			nameBytes = append(nameBytes, byte('.'))
			i += 2
			break
		}
		length := int(data[i])
		nameBytes = append(nameBytes, data[i+1:i+1+length]...)
		i += 1 + length
		if length == 0 {
			break
		}
		nameBytes = append(nameBytes, byte('.'))
	}
	nameBytes = nameBytes[:len(nameBytes)-1]
	return string(nameBytes), i
}

func QuestionFromBytes(data []byte, startIndex int) (Question, int) {
	q := Question{}
	i := startIndex
	q.Name, i = GetName(data, i)
	q.Type = uint16(data[i])<<8 | uint16(data[i+1])
	q.Class = uint16(data[i+2])<<8 | uint16(data[i+3])
	return q, i + 4
}

func (q Question) String() string {
	return fmt.Sprintf(
		"Question{name: %s, type: %d, class: %d}",
		q.Name,
		q.Type,
		q.Class,
	)
}
