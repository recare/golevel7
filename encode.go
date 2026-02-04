package golevel7

import (
	"errors"
	"io"
	"reflect"
)

// Encoder writes hl7 messages to a stream
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new Encoder that writes to stream w
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes the encoding of it to the stream
// It will panic if interface{} is not a pointer to a struct
func (e *Encoder) Encode(it interface{}) error {
	msg := &Message{}
	b, err := Marshal(msg, it)
	if err != nil {
		return err
	}
	i, err := e.w.Write(b)
	if err != nil {
		return err
	}
	if i < len(b) {
		return errors.New("Failed to write all bytes")
	}
	return nil
}

// Marshal will insert values into a message.
// It will panic if interface{} is not a pointer to a struct.
func Marshal(message *Message, it interface{}) ([]byte, error) {
	existingMSH, _ := message.Segment("MSH")

	// If we have no MSH header (in case of new message) add it first.
	if existingMSH == nil {
		segment := Segment{Value: []rune("MSH" + string(message.Delimeters.Field) + message.Delimeters.DelimeterField)}

		err := segment.parse(&message.Delimeters)
		if err != nil {
			return nil, err
		}

		message.Segments = append(message.Segments, segment)
	}

	baseStruct := reflect.ValueOf(it).Elem()

	baseStructType := baseStruct.Type()
	for i := 0; i < baseStruct.NumField(); i++ {
		fieldType := baseStructType.Field(i)

		fieldTag := fieldType.Tag.Get("hl7")
		if fieldTag == "" {
			continue
		}

		location := NewLocation(fieldTag)

		field := baseStruct.Field(i)

		switch field.Kind() {
		case reflect.String:
			if err := message.Set(location, field.String()); err != nil {
				return nil, err
			}
		}
	}

	return []byte(string(message.Value)), nil
}

func MarshalSegment(segment *Segment, it interface{}, delimeters *Delimeters) ([]byte, error) {
	baseStruct := reflect.ValueOf(it).Elem()

	baseStructType := baseStruct.Type()
	for i := 0; i < baseStruct.NumField(); i++ {
		fieldType := baseStructType.Field(i)

		fieldTag := fieldType.Tag.Get("hl7")
		if fieldTag == "" {
			continue
		}

		location := NewLocation(fieldTag)

		field := baseStruct.Field(i)

		switch field.Kind() {
		case reflect.String:
			if err := segment.SetForMarshaling(location, field.String(), delimeters); err != nil {
				return nil, err
			}
		}
	}

	return []byte(string(segment.Value)), nil
}

// SetForMarshaling will insert a value into a Segment at Location.
// ONLY use for MarshalSegment.
func (segment *Segment) SetForMarshaling(l *Location, val string, delimeters *Delimeters) error {
	if l.Segment == "" {
		return errors.New("Segment is required")
	}
	field, err := segment.Field(0)
	if err != nil || string(field.Value) != l.Segment {
		segment.forceField([]rune(l.Segment), 0)
	}

	segment.Set(l, val, delimeters)

	segment.Value = segment.encode(delimeters)
	return nil
}
