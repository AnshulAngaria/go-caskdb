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

func decodeRecord(data []byte) (uint32, string, string) {
	header := decodeHeader(data)

	key := string(data[headerSize : headerSize+header.KeySize])
	value := string(data[headerSize+header.KeySize : headerSize+header.KeySize+header.ValueSize])
	return header.Timestamp, key, value
}

func encodeRecord(timestamp uint32, key string, value string) (int, []byte) {
	header := encodeHeader(timestamp, uint32(len([]byte(key))), uint32(len([]byte(value))))
	data := append([]byte(key), []byte(value)...)
	return headerSize + len(data), append(header, data...)
}
