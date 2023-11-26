package main

type Packet struct {
	Header    Header
	Questions []Question
	Answers   []Answer
}

func PacketFromBytes(data []byte) Packet {
	p := Packet{}
	p.Header = HeaderFromBytes(data[:12])
	p.Questions = make([]Question, p.Header.QuestionCount)
	startIndex := 12
	for i := 0; i < int(p.Header.QuestionCount); i++ {
		p.Questions[i], startIndex = QuestionFromBytes(data, startIndex)
	}
	p.Answers = make([]Answer, p.Header.AnswerRecordCount)
	for i := 0; i < int(p.Header.AnswerRecordCount); i++ {
		p.Answers[i], startIndex = AnswerFromBytes(data, startIndex)
	}
	return p
}

func (p Packet) AsBytes() []byte {
	res := p.Header.AsBytes()
	for _, question := range p.Questions {
		res = append(res, question.AsBytes()...)
	}
	for _, answer := range p.Answers {
		res = append(res, answer.AsBytes()...)
	}
	return res
}
