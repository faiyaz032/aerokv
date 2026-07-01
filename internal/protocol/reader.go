package protocol

import (
	"bufio"
	"fmt"
)

func (p *Parser) readLine() (string, error) {
	line, err := p.reader.ReadSlice('\n')

	if err == bufio.ErrBufferFull {
		return "", fmt.Errorf("RESP line exceeds maximum length")
	}

	if err != nil {
		return "", err
	}

	if len(line) < 2 || line[len(line)-2] != '\r' {
		return "", fmt.Errorf("invalid RESP line: missing CRLF")
	}

	return string(line[:len(line)-2]), nil
}
