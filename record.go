package gocaskdb

type Record struct {
	Header Header
	Key    string
	Value  []byte
}

type Header struct {
	Checksum  uint32
	Timestamp uint32
	Expiry    uint32
	KeySize   uint32
	ValueSize uint32
}

type KeyDir map[string]Value

type Value struct {
	Timestamp  int
	RecordSize int
	RecordPos  int
	FileID     int
}
