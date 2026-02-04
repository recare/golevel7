package golevel7

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSegParse(t *testing.T) {
	val := []rune(`PID|||12001||Jones^John^^^Mr.||19670824|M|||123 West St.^^Denver^CO^80020^USA~520 51st Street^^Denver^CO^80020^USA|||||||
`)
	seps := NewDelimeters()
	seg := &Segment{Value: val}
	seg.parse(seps)
	if len(seg.Fields) != 20 {
		t.Errorf("Expected 20 fields got %d\n", len(seg.Fields))
	}
}

func TestSegSet(t *testing.T) {
	seps := NewDelimeters()
	loc := "ZZZ.10"
	l := NewLocation(loc)
	seg := &Segment{}
	err := seg.Set(l, "TEST", seps)
	if err != nil {
		t.Error(seg)
	}
	str, err := seg.Get(l)
	if err != nil {
		t.Error(err)
	}
	if str != "TEST" {
		t.Errorf("Expected TEST got %s\n", str)
	}
}

type obx struct {
	SegmentNumber             string `hl7:"OBX.1"`
	ValueType                 string `hl7:"OBX.2"`
	ObservationIdentifier     string `hl7:"OBX.3.0"`
	ObservationIdentifierText string `hl7:"OBX.3.1"`
	ObservationResultStatus   string `hl7:"OBX.11"`
	ObservationDate           string `hl7:"OBX.14"`
}

func TestMarshalSegment_BaseCase(t *testing.T) {
	testValue := []rune(`OBX|2|NM|RBC^RED BLOOD CELL COUNT||||||||F|||20050615230600`)

	newSegment := &Segment{}
	delimiters := NewDelimeters()

	newOBX := &obx{
		SegmentNumber:             "2",
		ValueType:                 "NM",
		ObservationIdentifier:     "RBC",
		ObservationIdentifierText: "RED BLOOD CELL COUNT",
		ObservationResultStatus:   "F",
		ObservationDate:           "20050615230600",
	}

	marshaledSegment, err := MarshalSegment(newSegment, newOBX, delimiters)
	require.NoError(t, err)

	assert.Equal(t, string(testValue), string(marshaledSegment))
}

func TestMarshalSegment_WorksLikeMarshal(t *testing.T) {
	delimeters := NewDelimeters()

	newMessage := &Message{Delimeters: *delimeters}

	newOBX := &obx{
		SegmentNumber:             "2",
		ValueType:                 "NM",
		ObservationIdentifier:     "RBC",
		ObservationIdentifierText: "RED BLOOD CELL COUNT",
		ObservationResultStatus:   "F",
		ObservationDate:           "20050615230600",
	}

	marshaledMessage, err := Marshal(newMessage, newOBX)
	require.NoError(t, err)

	newSegment := &Segment{}

	marshaledSegment, err := MarshalSegment(newSegment, newOBX, delimeters)
	require.NoError(t, err)

	marshaledMessageLines := bytes.Split(marshaledMessage, []byte{'\r'})
	assert.Equal(t, marshaledMessageLines[1], marshaledSegment)
}
