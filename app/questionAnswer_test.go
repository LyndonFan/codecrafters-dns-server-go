package main

import (
	"testing"
)

func TestGetNameBytes(t *testing.T) {
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
	actualName, actualIndex := GetNameBytes(data, 0)
	if actualName != expectedName {
		t.Errorf("Expected name to be %s, but got %s", expectedName, actualName)
	}
	if actualIndex != expectedIndex {
		t.Errorf("Expected index to be %d, but got %d", expectedIndex, actualIndex)
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
