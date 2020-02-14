package parser

import "testing"

func TestFrame(t *testing.T) {
	frame := []byte{}
	frame = inU8(frame, 5)
	frame, u := outU8_1(frame)
	if u != 5 {
		t.Error("failed!")
	} else {
		t.Error("good")
	}
}
