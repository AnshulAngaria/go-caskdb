package gocaskdb

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"
	"time"
)

type DiskStore struct {
	sync.Mutex
	file          *os.File
	writePosition int
	keyDir        KeyDir
}

func NewDiskStore(fileName string) (*DiskStore, error) {
	ds := &DiskStore{keyDir: make(KeyDir)}

	if isFileExists(fileName) {
		ds.initKeyDir(fileName)
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	ds.file = file
	return ds, nil
}

func (ds *DiskStore) initKeyDir(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	for {
		headerBytes := make([]byte, headerSize)
		_, err := io.ReadFull(file, headerBytes)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		header := decodeHeader(headerBytes)
		key := make([]byte, header.KeySize)
		value := make([]byte, header.ValueSize)
		_, err = io.ReadFull(file, key)

		if err != nil {
			return err
		}
		_, err = io.ReadFull(file, value)
		if err != nil {
			return err
		}

		totalSize := headerSize + header.KeySize + header.ValueSize
		ds.keyDir[string(key)] = NewValueEntry(header.Timestamp, uint32(ds.writePosition), totalSize)
		ds.writePosition += int(totalSize)
		fmt.Printf("loaded key=%s, value=%s\n", key, value)
	}
	return nil
}

func isFileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil || errors.Is(err, fs.ErrExist) {
		return true
	}
	return false
}

func (d *DiskStore) write(data []byte) {
	d.Lock()
	defer d.Unlock()
	if _, err := d.file.Write(data); err != nil {
		panic(err)
	}

	if err := d.file.Sync(); err != nil {
		panic(err)
	}
}

func (ds *DiskStore) Get(key string) string {
	ds.Lock()
	defer ds.Unlock()

	entry, ok := ds.keyDir[key]
	if !ok {
		return ""
	}

	ds.file.Seek(int64(entry.RecordPos), 0)
	data := make([]byte, entry.RecordSize)

	_, err := io.ReadFull(ds.file, data)
	if err != nil {
		panic("read error")
	}
	_, _, val := decodeRecord(data)
	return val

}

func (ds *DiskStore) Set(key, value string) {
	ds.Lock()
	defer ds.Unlock()

	timestamp := uint32(time.Now().Unix())
	size, data := encodeRecord(timestamp, key, value)
	ds.write(data)
	ds.keyDir[key] = NewValueEntry(timestamp, uint32(ds.writePosition), uint32(size))
	ds.writePosition += size
}

func (ds *DiskStore) Close() bool {
	ds.Lock()
	defer ds.Unlock()

	ds.file.Sync()
	if err := ds.file.Close(); err != nil {
		return false
	}
	return true
}
