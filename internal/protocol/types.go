package protocol

type Type int

const (
	SimpleString Type = iota
	Error
	Integer
	BulkString
	Array
)

type Value struct {
	Type    Type
	Text    string
	Integer int64
	Array   []Value
	Bulk    []byte
	Null    bool
}
