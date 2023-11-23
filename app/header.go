package main

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

func (h *Header) FromBytes(data []byte) {
	h.Identifier = uint16(data[0])<<8 | uint16(data[1])
	h.QR = data[2]&(1<<7) != 0
	h.OpCode = uint8(data[2]>>3) & 0x07
	h.Authoritative = data[2]&(1<<2) != 0
	h.Truncation = data[2]&(1<<1) != 0
	h.RecursionDesired = data[2]&(1) != 0
	h.Reserved = uint8(data[3] >> 4)
	h.ResponseCode = uint8(data[3] & 0x0f)
	h.QuestionCount = uint16(data[4])<<8 | uint16(data[5])
	h.AnswerRecordCount = uint16(data[6])<<8 | uint16(data[7])
	h.AuthorityRecordCount = uint16(data[8])<<8 | uint16(data[9])
	h.AdditionalRecordCount = uint16(data[10])<<8 | uint16(data[11])
	return
}
