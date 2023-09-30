package gocaskdb

import "encoding/binary"

type Header struct {
	//	Checksum  uint32
	Timestamp uint32
	//	Expiry    uint32
	KeySize   uint32
	ValueSize uint32
}

func encodeHeader(timestamp uint32, keySize uint32, valueSize uint32) []byte {
	header := make([]byte, headerSize)
	binary.LittleEndian.PutUint32(header[0:4], timestamp)
	binary.LittleEndian.PutUint32(header[4:8], keySize)
	binary.LittleEndian.PutUint32(header[8:12], valueSize)
	return header
}

func decodeHeader(header []byte) Header {
	timestamp := binary.LittleEndian.Uint32(header[0:4])
	keySize := binary.LittleEndian.Uint32(header[4:8])
	valueSize := binary.LittleEndian.Uint32(header[8:12])
	return Header{timestamp, keySize, valueSize}
}
