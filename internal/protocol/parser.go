package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(reader *bufio.Reader) *Parser {
	return &Parser{
		reader: reader,
	}
}

func (p *Parser) Parse() (Value, error) {
	prefix, err := p.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch prefix {
	case '+':
		return p.parseSimpleString()
	case '-':
		return p.parseError()
	case ':':
		return p.parseInteger()
	case '$':
		return p.parseBulkString()
	case '*':
		return p.parseArray()
	default:
		return Value{}, fmt.Errorf("unknown RESP type %q", prefix)
	}
}

func (p *Parser) parseSimpleString() (Value, error) {
	s, err := p.readLine()
	if err != nil {
		return Value{}, err
	}

	return Value{
		Type: SimpleString,
		Text: s,
	}, nil
}

func (p *Parser) parseError() (Value, error) {
	s, err := p.readLine()
	if err != nil {
		return Value{}, err
	}

	return Value{
		Type: Error,
		Text: s,
	}, nil
}

func (p *Parser) parseInteger() (Value, error) {
	s, err := p.readLine()
	if err != nil {
		return Value{}, err
	}

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return Value{}, err
	}

	return Value{
		Type:    Integer,
		Integer: n,
	}, nil
}

func (p *Parser) parseBulkString() (Value, error) {
	lengthStr, err := p.readLine()
	if err != nil {
		return Value{}, err
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return Value{}, err
	}

	// RESP null bulk string: $-1\r\n
	if length == -1 {
		return Value{
			Type: BulkString,
			Null: true,
		}, nil
	}

	if length < -1 {
		return Value{}, fmt.Errorf("invalid bulk length: %d", length)
	}
	if length > MaxBulkLength {
		return Value{}, fmt.Errorf(
			"bulk string too large: %d bytes (max %d)",
			length,
			MaxBulkLength,
		)
	}

	buf := make([]byte, length)

	_, err = io.ReadFull(p.reader, buf)
	if err != nil {
		return Value{}, err
	}

	// consume the trailing \r\n
	if _, err := p.reader.Discard(2); err != nil {
		return Value{}, err
	}

	return Value{
		Type: BulkString,
		Bulk: buf,
	}, nil
}

func (p *Parser) parseArray() (Value, error) {
	countStr, err := p.readLine()
	if err != nil {
		return Value{}, err
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return Value{}, err
	}

	// RESP null array: *-1\r\n
	if count == -1 {
		return Value{
			Type: Array,
			Null: true,
		}, nil
	}

	if count < -1 {
		return Value{}, fmt.Errorf("invalid array length: %d", count)
	}

	if count > MaxArrayLength {
		return Value{}, fmt.Errorf("array too large: %d elements (max %d)", count, MaxArrayLength)
	}

	values := make([]Value, count)

	for i := 0; i < count; i++ {
		value, err := p.Parse()
		if err != nil {
			return Value{}, err
		}

		values[i] = value
	}

	return Value{
		Type:  Array,
		Array: values,
	}, nil
}
