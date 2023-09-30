package gocaskdb

type Record struct {
	Header Header
	Key    string
	Value  []byte
}

type KeyDir map[string]ValueEntry

type ValueEntry struct {
	Timestamp  uint32
	RecordSize uint32
	RecordPos  uint32
}

func NewValueEntry(timestamp uint32, position uint32, totalSize uint32) ValueEntry {
	return ValueEntry{timestamp, position, totalSize}
}
