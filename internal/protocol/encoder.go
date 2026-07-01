package protocol

import (
	"bufio"
	"fmt"
	"io"
)

type Encoder struct {
	writer *bufio.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer: bufio.NewWriter(w),
	}
}

func (e *Encoder) Encode(value Value) error {
	switch value.Type {
	case SimpleString:
		return e.SimpleString(value.Text)

	case Error:
		return e.Error(value.Text)

	case Integer:
		return e.Integer(value.Integer)

	case BulkString:
		if value.Null {
			return e.NullBulkString()
		}
		return e.BulkString(value.Bulk)

	case Array:
		return e.Array(value.Array)

	default:
		return fmt.Errorf("unknown RESP type: %d", value.Type)
	}
}

func (e *Encoder) Flush() error {
	return e.writer.Flush()
}

func (e *Encoder) SimpleString(s string) error {
	_, err := fmt.Fprintf(e.writer, "+%s\r\n", s)
	return err
}

func (e *Encoder) Error(msg string) error {
	_, err := fmt.Fprintf(e.writer, "-%s\r\n", msg)
	return err
}

func (e *Encoder) Integer(n int64) error {
	_, err := fmt.Fprintf(e.writer, ":%d\r\n", n)
	return err
}

func (e *Encoder) BulkString(s []byte) error {
	if _, err := fmt.Fprintf(e.writer, "$%d\r\n", len(s)); err != nil {
		return err
	}

	if _, err := e.writer.Write(s); err != nil {
		return err
	}

	if _, err := e.writer.WriteString("\r\n"); err != nil {
		return err
	}

	return nil
}

func (e *Encoder) NullBulkString() error {
	_, err := e.writer.WriteString("$-1\r\n")
	return err
}

func (e *Encoder) Array(values []Value) error {
	_, err := fmt.Fprintf(e.writer, "*%d\r\n", len(values))
	if err != nil {
		return err
	}

	for _, value := range values {
		err := e.Encode(value)
		if err != nil {
			return err
		}
	}

	return nil
}
