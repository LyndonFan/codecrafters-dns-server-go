package main

import (
	"testing"
)

func TestGetName(t *testing.T) {
	data := []byte{
		0x03, 0x77, 0x77, 0x77,
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
		0x03, 0x63, 0x6f, 0x6d,
		0x00,
		0x00, 0x01,
		0x00, 0x01,
	}
	expectedName := "www.google.com"
	expectedIndex := 16
	actualName, actualIndex := GetName(data, 0)
	if actualName != expectedName {
		t.Errorf("Expected name to be %s, but got %s", expectedName, actualName)
	}
	if actualIndex != expectedIndex {
		t.Errorf("Expected index to be %d, but got %d", expectedIndex, actualIndex)
	}
}

func TestGetNameWithPointer(t *testing.T) {
	header := Header{
		Identifier:    0x04d2,
		QR:            false,
		OpCode:        0,
		QuestionCount: 4,
	}
	// example from https://www.rfc-editor.org/rfc/rfc1035#section-4.1.4
	questionData := []byte{
		// first question: F.ISI.ARPA
		1, 70, 3, 73, 83, 73, 4, 65, 82, 80, 65, 0, 0, 1, 0, 1,
		// second question: FOO.F.ISI.ARPA (with pointer)
		3, 70, 79, 79, 0xC0, 12, 0, 1, 0, 1,
		// third question: ARPA (with pointer)
		0xC0, 18, 0, 1, 0, 1,
		// fourth question: empty string / root
		0, 0, 1, 0, 1,
	}
	// concat header and questionData
	data := append(header.AsBytes(), questionData...)

	expectedQuestions := make([]Question, 4)
	expectedQuestions[0] = Question{
		Name:  "F.ISI.ARPA",
		Type:  1,
		Class: 1,
	}
	expectedQuestions[1] = Question{
		Name:  "FOO.F.ISI.ARPA",
		Type:  1,
		Class: 1,
	}
	expectedQuestions[2] = Question{
		Name:  "ARPA",
		Type:  1,
		Class: 1,
	}
	expectedQuestions[3] = Question{
		Name:  "",
		Type:  1,
		Class: 1,
	}

	actualQuestions := make([]Question, 4)
	startIndex := 12
	for i := 0; i < 4; i++ {
		actualQuestions[i], startIndex = QuestionFromBytes(data, startIndex)
	}
	if startIndex != len(data) {
		t.Errorf("Expected index to be %d, but got %d", len(data), startIndex)
	}

	for i := 0; i < 4; i++ {
		if actualQuestions[i].Name != expectedQuestions[i].Name {
			t.Errorf("Expected name of question %d to be %s, but got %s", i, expectedQuestions[i].Name, actualQuestions[i].Name)
		}
		if actualQuestions[i].Type != expectedQuestions[i].Type {
			t.Errorf("Expected type of question %d to be %d, but got %d", i, expectedQuestions[i].Type, actualQuestions[i].Type)
		}
		if actualQuestions[i].Class != expectedQuestions[i].Class {
			t.Errorf("Expected class of question %d to be %d, but got %d", i, expectedQuestions[i].Class, actualQuestions[i].Class)
		}
	}
}

func TestQuestionFromBytes(t *testing.T) {
	data := []byte{
		0x03, 0x77, 0x77, 0x77,
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
		0x03, 0x63, 0x6f, 0x6d,
		0x00,
		0x00, 0x01,
		0x00, 0x01,
	}

	expectedQuestion := Question{
		Name:  "www.google.com",
		Type:  1,
		Class: 1,
	}

	question, _ := QuestionFromBytes(data, 0)

	if question.Name != expectedQuestion.Name {
		t.Errorf("Expected name to be %s, but got %s", expectedQuestion.Name, question.Name)
	}

	if question.Type != expectedQuestion.Type {
		t.Errorf("Expected type to be %d, but got %d", expectedQuestion.Type, question.Type)
	}

	if question.Class != expectedQuestion.Class {
		t.Errorf("Expected class to be %d, but got %d", expectedQuestion.Class, question.Class)
	}
}
